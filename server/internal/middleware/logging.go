package middleware

import (
	"log"
	"net/http"
	"time"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		route := r.URL.Path
		log.Printf("request_id=%s method=%s route=%s status=%d duration_ms=%d",
			RequestIDFromContext(r.Context()),
			r.Method,
			route,
			rec.status,
			time.Since(start).Milliseconds(),
		)
	})
}
