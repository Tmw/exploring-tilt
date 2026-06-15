package middlewares

import (
	"net/http"
	"strings"

	"github.com/tmw/exploring-tilt/pkg/middleware"
)

// basic CORS middleware
type CorsConfig struct {
	Origin  string
	Methods []string
	Headers []string
}

func (c CorsConfig) MethodsList() string {
	return strings.Join(c.Methods, ", ")
}

func (c CorsConfig) HeadersList() string {
	return strings.Join(c.Headers, ", ")
}

func Cors(config CorsConfig) middleware.Middleware {
	return middleware.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", config.Origin)
			w.Header().Set("Access-Control-Allow-Methods", config.MethodsList())
			w.Header().Set("Access-Control-Allow-Headers", config.HeadersList())

			// handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}
