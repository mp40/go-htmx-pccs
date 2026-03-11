package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type Data interface {
	CountUserByID(ID uuid.UUID) (int, error)
}

type Middleware struct {
	data Data
}

func NewMiddlewareService(data Data) *Middleware {
	return &Middleware{
		data: data,
	}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		ctx := context.WithValue(c, "Test", "123")

		cookie, err := r.Cookie("userID")
		if err != nil {
			// not autherised || internal server error
		}
		if cookie == nil {
			// not signed in or authorised
		}
		if cookie != nil {
			rawID := cookie.Value
			userID, err := uuid.Parse(rawID)
			if err != nil {
				// not autherised || internal server error
			}

			// get count by id
			count, err := m.data.CountUserByID(userID)
			if err != nil {
				// not autherised || internal server error
			}
			if count > 1 {
				// ID is unique
				// not autherised || internal server error
			}
			if count == 0 {
				// we set cookie so if no matching ID?
				// not autherised || internal server erro
			}

			// if user count = 1 authorised and signed in
		}

		next.ServeHTTP(w, r.WithContext(ctx))
		// return
	})
}
