package middleware

import (
	"context"
	"net/http"
)

type Middleware struct{}

func NewMiddlewareService() *Middleware {
	return &Middleware{}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		ctx := context.WithValue(c, "Test", "123")
		next.ServeHTTP(w, r.WithContext(ctx))
		// return
	})
}
