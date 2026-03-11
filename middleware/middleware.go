package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const userContextKey contextKey = "userID"

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
		cookie, err := r.Cookie(string(userContextKey))
		if err != nil {
			// GREY - for now fall through
			next.ServeHTTP(w, r)
			return
		}
		rawID := cookie.Value
		userID, err := uuid.Parse(rawID)
		if err != nil {
			// GREY if we set cookie, then we expect parseable id
			// for now fall through - latter decide other action ie redirect to sign in with err msg
			next.ServeHTTP(w, r)
			return
		}

		// get count by id
		count, err := m.data.CountUserByID(userID)
		if err != nil {
			// GREY data layer error - do not know correct count
			// for now fall through
			next.ServeHTTP(w, r)
			return
		}
		if count > 1 {
			// GREY if we set cookie, then we expect 1
			// for now fall through - latter decide other action ie redirect to sign in
			next.ServeHTTP(w, r)
			return
		}
		if count == 0 {
			// GREY if we set cookie, then we expect 1
			// for now fall through
			next.ServeHTTP(w, r)
			return
		}

		// if user count = 1, def YES
		c := r.Context()
		ctx := context.WithValue(c, userContextKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// should these and the key be in Auth? or somewhere else? Or through interface or dep inj
func GetUserID(ctx context.Context) *uuid.UUID {
	userID, ok := ctx.Value(userContextKey).(uuid.UUID)
	if !ok {
		return nil
	}
	return &userID
}

func IsSignedIn(ctx context.Context) bool {
	_, ok := ctx.Value(userContextKey).(uuid.UUID)
	if !ok {
		return false
	}
	return true
}
