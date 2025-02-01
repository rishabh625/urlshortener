package middleware

import (
	"log"
	"net/http"
	"time"
)

func LatencyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		duration := time.Since(startTime)
		log.Printf("Method: %s, Path: %s, ROUTE pattern: %s, Latency: %v",
			r.Method, r.URL.Path, r.Pattern, duration)
		h.ServeHTTP(w, r)
	})
}
