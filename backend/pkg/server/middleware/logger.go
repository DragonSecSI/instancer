package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

func Logger(l *zerolog.Logger) func(next http.Handler) http.Handler {
	logger := l.With().Str("module", "server").Logger()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			r = r.WithContext(logger.WithContext(r.Context()))

			rw := NewResponseWriter(w)

			next.ServeHTTP(rw, r)

			panicVal := recover()
			if panicVal != nil {
				rw.statusCode = http.StatusInternalServerError
				defer panic(panicVal)
			}

			logger.
				Info().
				Str("method", r.Method).
				Str("url", r.URL.RequestURI()).
				Str("user_agent", r.UserAgent()).
				Int("status_code", rw.statusCode).
				Dur("elapsed_time", time.Since(start)).
				Msg("Request")
		})
	}
}
