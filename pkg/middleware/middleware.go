package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

type Stack struct {
	middlewares []Middleware
}

func New(stack ...Middleware) *Stack {
	return &Stack{
		middlewares: stack,
	}
}

func (s *Stack) Wrap(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := http.Handler(handler)
		for idx := range len(s.middlewares) {
			next := s.middlewares[len(s.middlewares)-idx-1]
			wrapped = next(wrapped)
		}

		wrapped.ServeHTTP(w, r)
	})
}
