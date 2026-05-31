package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/session"
)

type contextKey string

const userContextKey contextKey = "userID"
const sessionCookieKey contextKey = "sessionID"

type Session interface {
	// down the road make Session DTO to replace session.Session - and other structs in code base, will help decouple
	GetValidSessionByID(ID uuid.UUID) (*session.Session, error)
}

type Middleware struct {
	session Session
}

func NewMiddlewareService(session Session) *Middleware {
	return &Middleware{
		session: session,
	}
}

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(string(sessionCookieKey))
		if err != nil {
			// GREY - for now fall through
			next.ServeHTTP(w, r)
			return
		}
		rawID := cookie.Value
		sessionID, err := uuid.Parse(rawID)
		if err != nil {
			slog.Warn("middleware, unparseable session id", "rawID", rawID)
			// GREY if we set cookie, then we expect parseable id
			// for now fall through - latter decide other action ie redirect to sign in with err msg
			next.ServeHTTP(w, r)
			return
		}
		session, err := m.session.GetValidSessionByID(sessionID)
		if err != nil {
			slog.Error("middleware, session data error", "err", err)
			// GREY data layer error
			// for now fall through
			next.ServeHTTP(w, r)
			return
		}
		if session == nil {
			slog.Warn("middleware, session not found")
			// for now fall through
			next.ServeHTTP(w, r)
			return
		}

		c := r.Context()
		ctx := context.WithValue(c, userContextKey, session.UserID)

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

func GetSessionID(ctx context.Context) *uuid.UUID {
	sessionID, ok := ctx.Value(sessionCookieKey).(uuid.UUID)
	if !ok {
		return nil
	}
	return &sessionID
}

func IsSignedIn(ctx context.Context) bool {
	_, ok := ctx.Value(userContextKey).(uuid.UUID)
	if !ok {
		return false
	}
	return true
}
