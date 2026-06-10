package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type renderSignUpModal interface {
	RenderSignUpModal(w io.Writer) error
}

func getSignUpHandler(render renderSignUpModal) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		if isHtmx {
			err := render.RenderSignUpModal(w)
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

		w.Header().Set("Content-Type", "text/html")
	})
}

type postSignUpRender interface {
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type postSignUpAuth interface {
	SignUp(email string, password string) (*uuid.UUID, error)
}

type postSignUpSession interface {
	AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error)
}

func postSignUpHandler(render postSignUpRender, auth postSignUpAuth, session postSignUpSession) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))

		_, err := mail.ParseAddress(email)
		if err != nil || len(password) < 8 {
			err = render.RenderErrorMessageFragment(w, "invalid sign up: provide email and password at least 8 characters long")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			return
		}

		newID, err := auth.SignUp(email, password)
		if err != nil {
			err = render.RenderErrorMessageFragment(w, "nfi")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			return
		}

		if newID == nil {
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			return
		}
		// if yes
		now := time.Now()
		sessionID, err := session.AddSession(*newID, now.Add(12*time.Hour))
		if err != nil {
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			return
		}
		cookie := getSecureCookie(sessionID)
		http.SetCookie(w, &cookie)
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
	})
}
