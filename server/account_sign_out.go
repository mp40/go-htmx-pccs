package server

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type deleteSignOutSession interface {
	DeleteSessionByID(ID uuid.UUID) error
}

func deleteSignOutHandler(session deleteSignOutSession) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			redirectWithExpiredCookie(w)
			return
		}

		err = session.DeleteSessionByID(sessionID)
		if err != nil {
			slog.Error("delete session error", "err", err)
		}

		redirectWithExpiredCookie(w)
	})
}

func redirectWithExpiredCookie(w http.ResponseWriter) {
	expiredCookie := removeSecureCookie()
	http.SetCookie(w, &expiredCookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

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
