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

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(string(sessionCookieKey))
		if err != nil {
			slog.Error("middleware, get cookie error", "err", err)
			next.ServeHTTP(w, r)
			return
		}
		rawID := cookie.Value
		sessionID, err := uuid.Parse(rawID)
		if err != nil {
			slog.Warn("middleware, unparseable session id", "rawID", rawID)
			expired := removeSecureCookie()
			http.SetCookie(w, &expired)
			next.ServeHTTP(w, r)
			return
		}
		session, err := m.session.GetValidSessionByID(sessionID)
		if err != nil {
			slog.Error("middleware, session data error", "err", err)
			next.ServeHTTP(w, r)
			return
		}
		if session == nil {
			slog.Warn("middleware, session not found")
			next.ServeHTTP(w, r)
			return
		}

		c := r.Context()
		ctx := m.enrichContextFunc(c, session.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DUP - removeSecureCookie
func removeSecureCookie() http.Cookie {
	cookie := http.Cookie{
		Name:     "sessionID",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
}
