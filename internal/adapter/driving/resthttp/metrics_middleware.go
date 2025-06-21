package resthttp

import (
	"net/http"
	"strconv"
	"time"

	"ghorkov32/proletariat-budget-be/internal/common"
	"github.com/go-chi/chi/v5"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsCollector(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rw := NewResponseWriter(w)

			begin := time.Now()
			next.ServeHTTP(rw, r)
			elapsed := time.Since(begin)

			path := chi.RouteContext(r.Context()).RoutePattern()
			common.HTTPDuration.WithLabelValues(path).Observe(elapsed.Seconds())
			common.RequestCounter.WithLabelValues(path, strconv.Itoa(rw.statusCode)).Inc()
		})
}
