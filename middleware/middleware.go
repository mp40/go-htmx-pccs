package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/session"
)

type contextKey string

const sessionCookieKey contextKey = "sessionID"

type EnrichContextWithUserIDFunc func(ctx context.Context, userID uuid.UUID) context.Context

type Session interface {
	// down the road make Session DTO to replace session.Session - and other structs in code base, will help decouple
	GetValidSessionByID(ID uuid.UUID) (*session.Session, error)
}

type Middleware struct {
	session           Session
	enrichContextFunc EnrichContextWithUserIDFunc
}

func NewMiddlewareService(session Session, cb EnrichContextWithUserIDFunc) *Middleware {
	return &Middleware{
		session:           session,
		enrichContextFunc: cb,
	}
}

// TODO
// Think about when we have public and private pages
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
		ctx := m.enrichContextFunc(c, session.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSessionID(ctx context.Context) *uuid.UUID {
	sessionID, ok := ctx.Value(sessionCookieKey).(uuid.UUID)
	if !ok {
		return nil
	}
	return &sessionID
}
