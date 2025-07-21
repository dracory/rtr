package middlewares

import (
	"github.com/dracory/rtr"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// CompressMiddleware returns a middleware that compresses HTTP responses.
// It supports gzip, deflate, and brotli compression based on the client's Accept-Encoding header.
// This is a thin wrapper around klauspost/compress/gzhttp's Transport.
func CompressMiddleware(level int, types ...string) rtr.MiddlewareInterface {
	return rtr.NewMiddleware().
		SetName("Compress").
		SetHandler(chimiddleware.Compress(level, types...))
}
