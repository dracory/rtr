package middlewares

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/dracory/rtr"
)

// defaultCompressibleContentTypes lists the content types that are compressed
// by default when the caller does not specify a custom list.
var defaultCompressibleContentTypes = []string{
	"text/html",
	"text/css",
	"text/plain",
	"text/javascript",
	"application/javascript",
	"application/x-javascript",
	"application/json",
	"application/atom+xml",
	"application/rss+xml",
	"application/xml",
	"text/xml",
	"image/svg+xml",
}

// CompressMiddleware returns a middleware that compresses HTTP responses.
// It supports gzip and deflate compression based on the client's Accept-Encoding header.
// This is a stdlib-only reimplementation of chi's middleware.Compress.
func CompressMiddleware(level int, types ...string) rtr.MiddlewareInterface {
	compressor := newCompressor(level, types...)
	return rtr.NewMiddleware().
		SetName("Compress").
		SetHandler(compressor.handler)
}

// encoderFunc wraps an io.Writer with a streaming compression algorithm.
type encoderFunc func(w io.Writer, level int) io.Writer

// ioResetterWriter is an io.Writer that can be Reset to a new underlying writer.
type ioResetterWriter interface {
	io.Writer
	Reset(w io.Writer)
}

// compressor holds the encoding configuration for the compress middleware.
type compressor struct {
	encoders           map[string]encoderFunc
	pooledEncoders     map[string]*sync.Pool
	allowedTypes       map[string]struct{}
	allowedWildcards   map[string]struct{}
	encodingPrecedence []string
	level              int
}

// newCompressor creates a new compressor with the given compression level and
// optional content-type filter. If no types are given, a default list is used.
// It panics if the compression level is invalid for gzip or deflate.
func newCompressor(level int, types ...string) *compressor {
	// Validate the level early so we fail fast instead of producing nil encoders
	// that would cause silent data corruption (Content-Encoding set without compression).
	if err := validateCompressionLevel(level); err != nil {
		panic(err.Error())
	}

	allowedTypes := make(map[string]struct{})
	allowedWildcards := make(map[string]struct{})

	if len(types) > 0 {
		for _, t := range types {
			if strings.Contains(strings.TrimSuffix(t, "/*"), "*") {
				panic("middleware/compress: Unsupported content-type wildcard pattern. Only '/*' supported")
			}
			if before, ok := strings.CutSuffix(t, "/*"); ok {
				allowedWildcards[before] = struct{}{}
			} else {
				allowedTypes[t] = struct{}{}
			}
		}
	} else {
		for _, t := range defaultCompressibleContentTypes {
			allowedTypes[t] = struct{}{}
		}
	}

	c := &compressor{
		level:            level,
		encoders:         make(map[string]encoderFunc),
		pooledEncoders:   make(map[string]*sync.Pool),
		allowedTypes:     allowedTypes,
		allowedWildcards: allowedWildcards,
	}

	// Add deflate first, then gzip, so gzip has higher precedence.
	c.setEncoder("deflate", encoderDeflate)
	c.setEncoder("gzip", encoderGzip)

	return c
}

// setEncoder registers an encoder for a given encoding name.
func (c *compressor) setEncoder(encoding string, fn encoderFunc) {
	encoding = strings.ToLower(encoding)
	if encoding == "" {
		panic("the encoding can not be empty")
	}
	if fn == nil {
		panic("attempted to set a nil encoder function")
	}

	delete(c.pooledEncoders, encoding)
	delete(c.encoders, encoding)

	encoder := fn(io.Discard, c.level)
	if _, ok := encoder.(ioResetterWriter); ok {
		pool := &sync.Pool{
			New: func() interface{} {
				return fn(io.Discard, c.level)
			},
		}
		c.pooledEncoders[encoding] = pool
	}
	if _, ok := c.pooledEncoders[encoding]; !ok {
		c.encoders[encoding] = fn
	}

	for i, v := range c.encodingPrecedence {
		if v == encoding {
			c.encodingPrecedence = append(c.encodingPrecedence[:i], c.encodingPrecedence[i+1:]...)
			break
		}
	}

	c.encodingPrecedence = append([]string{encoding}, c.encodingPrecedence...)
}

// handler returns the middleware handler that compresses responses.
func (c *compressor) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder, encoding, cleanup := c.selectEncoder(r.Header, w)

		cw := &compressResponseWriter{
			ResponseWriter:   w,
			w:                w,
			contentTypes:     c.allowedTypes,
			contentWildcards: c.allowedWildcards,
			encoding:         encoding,
			compressible:     false,
		}
		if encoder != nil {
			cw.w = encoder
		}
		defer cleanup()
		defer cw.close()

		next.ServeHTTP(cw, r)
	})
}

// selectEncoder picks the best encoder based on the Accept-Encoding header.
func (c *compressor) selectEncoder(h http.Header, w io.Writer) (io.Writer, string, func()) {
	header := h.Get("Accept-Encoding")
	accepted := strings.Split(strings.ToLower(header), ",")

	for _, name := range c.encodingPrecedence {
		if matchAcceptEncoding(accepted, name) {
			if pool, ok := c.pooledEncoders[name]; ok {
				encoder := pool.Get().(ioResetterWriter)
				cleanup := func() { pool.Put(encoder) }
				encoder.Reset(w)
				return encoder, name, cleanup
			}
			if fn, ok := c.encoders[name]; ok {
				return fn(w, c.level), name, func() {}
			}
		}
	}

	return nil, "", func() {}
}

// validateCompressionLevel checks that the level is valid for both gzip and flate.
// gzip accepts levels -2 (HuffmanOnly) to 9 (BestCompression), plus -1 (DefaultCompression).
// flate accepts levels 0 (NoCompression) to 9 (BestCompression), plus -1 (DefaultCompression).
// The intersection used by both is: -2, -1, 0..9.
func validateCompressionLevel(level int) error {
	switch {
	case level == gzip.DefaultCompression: // -1
		return nil
	case level == gzip.HuffmanOnly: // -2
		return nil
	case level >= gzip.NoCompression && level <= gzip.BestCompression: // 0..9
		return nil
	default:
		return fmt.Errorf("compress: invalid compression level %d (valid range: -2 to 9)", level)
	}
}

// matchAcceptEncoding checks if the encoding is in the accepted list.
// It properly parses Accept-Encoding tokens per RFC 7231:
//   - Each token is split on ";" to separate the encoding from q-values.
//   - Whitespace is trimmed.
//   - q=0 means "explicitly not acceptable" and is rejected.
//   - Substring matches are avoided (e.g. "xgzip" does not match "gzip").
func matchAcceptEncoding(accepted []string, encoding string) bool {
	for _, v := range accepted {
		// Trim whitespace
		v = strings.TrimSpace(v)
		// Split on ";" to separate encoding from parameters like q-values
		token, rest, _ := strings.Cut(v, ";")
		token = strings.TrimSpace(token)
		if token != encoding {
			continue
		}
		// Check q-value: q=0 means "not acceptable"
		if rest != "" {
			rest = strings.TrimSpace(rest)
			if strings.HasPrefix(rest, "q=") {
				qStr := strings.TrimPrefix(rest, "q=")
				qStr = strings.TrimSpace(qStr)
				if qStr == "0" || strings.HasPrefix(qStr, "0.") {
					// q=0 or q=0.x means explicitly not acceptable
					return false
				}
			}
		}
		return true
	}
	return false
}

// compressResponseWriter wraps http.ResponseWriter to compress the response
// when the content type is compressible.
type compressResponseWriter struct {
	http.ResponseWriter
	w                io.Writer
	contentTypes     map[string]struct{}
	contentWildcards map[string]struct{}
	encoding         string
	wroteHeader      bool
	compressible     bool
}

// isCompressible checks if the response's Content-Type is in the allowed list.
// Content types are matched case-insensitively per RFC 7231, and whitespace
// around the type is trimmed.
func (cw *compressResponseWriter) isCompressible() bool {
	contentType := cw.Header().Get("Content-Type")
	contentType, _, _ = strings.Cut(contentType, ";")
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	if _, ok := cw.contentTypes[contentType]; ok {
		return true
	}
	if contentType, _, hadSlash := strings.Cut(contentType, "/"); hadSlash {
		_, ok := cw.contentWildcards[contentType]
		return ok
	}
	return false
}

// WriteHeader sets the Content-Encoding header if the response is compressible.
func (cw *compressResponseWriter) WriteHeader(code int) {
	if cw.wroteHeader {
		cw.ResponseWriter.WriteHeader(code)
		return
	}
	cw.wroteHeader = true
	defer cw.ResponseWriter.WriteHeader(code)

	if cw.Header().Get("Content-Encoding") != "" {
		return
	}

	if !cw.isCompressible() {
		cw.compressible = false
		return
	}

	if cw.encoding != "" {
		cw.compressible = true
		cw.Header().Set("Content-Encoding", cw.encoding)
		cw.Header().Add("Vary", "Accept-Encoding")
		cw.Header().Del("Content-Length")
	}
}

// Write writes the data to the underlying writer, compressing if applicable.
func (cw *compressResponseWriter) Write(p []byte) (int, error) {
	if !cw.wroteHeader {
		cw.WriteHeader(http.StatusOK)
	}
	return cw.writer().Write(p)
}

// writer returns the active writer (compressor or raw response writer).
func (cw *compressResponseWriter) writer() io.Writer {
	if cw.compressible {
		return cw.w
	}
	return cw.ResponseWriter
}

// Flush flushes the underlying writers if they support flushing.
func (cw *compressResponseWriter) Flush() {
	if f, ok := cw.writer().(http.Flusher); ok {
		f.Flush()
	}
	if f, ok := cw.writer().(compressFlusher); ok {
		f.Flush()
		if f, ok := cw.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// Hijack supports connection hijacking if the underlying writer supports it.
func (cw *compressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := cw.writer().(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("compress: http.Hijacker is unavailable on the writer")
}

// Push supports HTTP/2 server push if the underlying writer supports it.
func (cw *compressResponseWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := cw.writer().(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return errors.New("compress: http.Pusher is unavailable on the writer")
}

// close closes the underlying encoder if it is a WriteCloser.
// Returns nil if the writer is not a WriteCloser (e.g. the raw ResponseWriter
// for non-compressible responses).
func (cw *compressResponseWriter) close() error {
	if c, ok := cw.writer().(io.WriteCloser); ok {
		return c.Close()
	}
	return nil
}

// Unwrap returns the underlying ResponseWriter for compatibility.
func (cw *compressResponseWriter) Unwrap() http.ResponseWriter {
	return cw.ResponseWriter
}

// compressFlusher is implemented by encoders that support flushing.
type compressFlusher interface {
	Flush() error
}

// encoderGzip creates a gzip writer at the given level.
func encoderGzip(w io.Writer, level int) io.Writer {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil
	}
	return gw
}

// encoderDeflate creates a deflate writer at the given level.
func encoderDeflate(w io.Writer, level int) io.Writer {
	dw, err := flate.NewWriter(w, level)
	if err != nil {
		return nil
	}
	return dw
}
