package httprepo

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func RequestLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			defer func() {
				duration := time.Since(start)

				logger.Info().Str("method", r.Method).Str("Path", r.URL.Path).Dur("duration", duration).Msg("request completed")
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func EnableCORS() func(http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
	
			next.ServeHTTP(w, r)
		})
	}
}

