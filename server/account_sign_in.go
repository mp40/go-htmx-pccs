package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/store"
)

type renderSignInModal interface {
	RenderSignInModal(w io.Writer) error
}

func getSignInHandler(render renderSignInModal) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderSignInModal(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// log err
			// ???
			// what to do if not htmx?
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

type postSignInRender interface {
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type postSignInAuth interface {
	SignIn(email string, password string) (*store.User, error)
}

type postSignInSession interface {
	AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error)
}

func postSignInHandler(render postSignInRender, auth postSignInAuth, session postSignInSession) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		user, err := auth.SignIn(email, password)
		if err != nil {
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		if user == nil {
			err = render.RenderErrorMessageFragment(w, "invalid sign in: check email and password")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		// if yes
		now := time.Now()
		sessionID, err := session.AddSession(user.ID, now.Add(12*time.Hour))
		if err != nil {
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		cookie := getSecureCookie(sessionID)
		http.SetCookie(w, &cookie)
		w.Header().Set("HX-Redirect", "/account")
		w.WriteHeader(http.StatusSeeOther)
	})
}
