package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &WriterWrapper{ResponseWriter: w, StatusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"ip":         getClientIP(r),
			"user_agent": r.UserAgent(),
			"status":     wrapped.StatusCode,
			"duration":   duration.Milliseconds(),
			"timestamp":  start.Format(time.RFC3339),
		}).Info("HTTP request")
	})
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}
