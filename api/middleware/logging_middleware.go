package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func LoggingMiddleware(log logger.AppLoggerInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.ErrorWithFields(r.Context(), "error!!!", err.(error), logger.FieldsSet{
						"err":   err,
						"trace": debug.Stack(),
					})
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)

			log.InfoWithFields(r.Context(), "Incoming HTTP request", logger.FieldsSet{
				"component": "request/start",
				"method":    r.Method,
				"path":      r.RequestURI,
				"ip":        r.RemoteAddr,
				"agent":     r.UserAgent(),
			})

			next.ServeHTTP(wrapped, r)

			log.InfoWithFields(r.Context(), "Finished handling HTTP request", logger.FieldsSet{
				"component": "request/end",
				"status":    wrapped.status,
				"method":    r.Method,
				"path":      r.URL.EscapedPath(),
				"duration":  time.Since(start),
				"ip":        r.RemoteAddr,
			})
		}

		return http.HandlerFunc(fn)
	}
}
