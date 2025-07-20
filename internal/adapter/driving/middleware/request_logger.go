package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// getLogLevelForStatus returns the appropriate log level based on HTTP status code
func getLogLevelForStatus(status int) zerolog.Level {
	switch {
	case status >= http.StatusInternalServerError:
		return zerolog.ErrorLevel
	case status >= http.StatusBadRequest:
		return zerolog.WarnLevel
	case status >= http.StatusMultipleChoices:
		return zerolog.InfoLevel
	case status >= http.StatusOK:
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}

// getLogMessageForStatus returns an appropriate message based on HTTP status code
func getLogMessageForStatus(status int) string {
	switch {
	case status >= http.StatusInternalServerError:
		return "HTTP request completed with server error"
	case status >= http.StatusBadRequest:
		return "HTTP request completed with client error"
	case status >= http.StatusMultipleChoices:
		return "HTTP request completed with redirect"
	case status >= http.StatusOK:
		return "HTTP request completed successfully"
	default:
		return "HTTP request completed with informational response"
	}
}

// DetailedRequestLogger provides more detailed logging including request headers
func DetailedRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(
				w,
				r.ProtoMajor,
			)

			// Log detailed request information
			logger := log.With().
				Str(
					"method",
					r.Method,
				).
				Str(
					"url",
					r.URL.String(),
				).
				Str(
					"remote_addr",
					r.RemoteAddr,
				).
				Str(
					"user_agent",
					r.UserAgent(),
				).
				Str(
					"request_id",
					middleware.GetReqID(r.Context()),
				).
				Str(
					"proto",
					r.Proto,
				).
				Int64(
					"content_length",
					r.ContentLength,
				).
				Logger()

			// Log request headers (be careful with sensitive data)
			headers := make(map[string]string)
			for name, values := range r.Header {
				if len(values) > 0 {
					// Skip sensitive headers
					if name != "Authorization" && name != "Cookie" {
						headers[name] = values[0]
					}
				}
			}

			logger.Info().
				Interface(
					"headers",
					headers,
				).
				Msg("HTTP request started")

			next.ServeHTTP(
				ww,
				r,
			)

			duration := time.Since(start)
			status := ww.Status()

			// Log with appropriate level based on status code
			logLevel := getLogLevelForStatus(status)
			message := getLogMessageForStatus(status)

			logger.WithLevel(logLevel).
				Int(
					"status",
					status,
				).
				Int(
					"bytes",
					ww.BytesWritten(),
				).
				Dur(
					"duration",
					duration,
				).
				Msg(message)
		},
	)
}
