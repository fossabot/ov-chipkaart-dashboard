package middlewares

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-kit/kit/log"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// Status returns the status of the response writer
func (rw *responseWriter) Status() int {
	return rw.status
}

// WhiteHeader writes the header
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			_ = logger.Log(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
