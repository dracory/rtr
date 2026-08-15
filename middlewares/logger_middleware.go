package middlewares

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dracory/rtr"
)

// LoggerMiddleware returns a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was, and how long it took.
// This is a stdlib-only reimplementation of chi's middleware.Logger.
func LoggerMiddleware() rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Logger").
		SetHandler(loggerHandler)
}

// loggerHandler is the stdlib-only implementation of the request logger
// previously provided by chi's middleware.Logger.
func loggerHandler(next http.Handler) http.Handler {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := newWrapResponseWriter(w)

		t1 := time.Now()

		// Log the request start
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		logger.Printf("\"%s %s://%s%s %s\" from %s - ",
			r.Method, scheme, r.Host, r.RequestURI, r.Proto, r.RemoteAddr)

		next.ServeHTTP(ww, r)

		// Log the request completion
		logger.Printf("%03d %dB in %s", ww.status, ww.bytesWritten, time.Since(t1))
	})
}

// wrapResponseWriter captures the HTTP status code and bytes written.
// It delegates http.Hijacker, http.Pusher, and http.Flusher to the underlying
// ResponseWriter so that downstream handlers (e.g. WebSocket upgraders) work
// correctly through the logger middleware.
type wrapResponseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int
	wroteHeader  bool
}

// newWrapResponseWriter creates a new wrapResponseWriter.
func newWrapResponseWriter(w http.ResponseWriter) *wrapResponseWriter {
	return &wrapResponseWriter{ResponseWriter: w}
}

// WriteHeader captures the status code and delegates to the embedded ResponseWriter.
func (w *wrapResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		w.ResponseWriter.WriteHeader(code)
		return
	}
	w.wroteHeader = true
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written and delegates to the embedded ResponseWriter.
func (w *wrapResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytesWritten += n
	return n, err
}

// Flush supports http.Flusher if the underlying writer implements it.
func (w *wrapResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack supports connection hijacking if the underlying writer implements it.
// This is required for WebSocket upgrades and similar use cases.
func (w *wrapResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Push supports HTTP/2 server push if the underlying writer implements it.
func (w *wrapResponseWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := w.ResponseWriter.(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Unwrap returns the underlying ResponseWriter for compatibility.
func (w *wrapResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
