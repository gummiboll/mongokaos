package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gummiboll/mongokaos/state"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &StatusRecorder{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)
		log.Printf(`%s %s %s %s [%d]`, r.Method, r.URL, r.Method, time.Since(start).Truncate(time.Millisecond), recorder.StatusCode)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-key") != state.GetAppState().Config.APIKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
