package middlewares

import "net/http"

func RequestSizeLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set a maximum body size of 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 1e+7)
		next.ServeHTTP(w, r)
	})
}
