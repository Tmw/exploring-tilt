package middlewares

import (
	"net/http"

	"github.com/tmw/exploring-tilt/pkg/middleware"
)

func JSON() middleware.Middleware {
	return middleware.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
}
