package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/data"
	"github.com/mp40/go-htmx-pccs/middleware"
)

type Render interface {
	RenderHomePage(w io.Writer, signedIn bool) error
	RenderHomeFragment(w io.Writer) error
	RenderAccountPage(w io.Writer) error
	RenderAccountFragment(w io.Writer) error
	RenderSignInModal(w io.Writer) error
	RenderSignUpModal(w io.Writer) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type Auth interface {
	SignIn(email string, password string) (*data.User, error)
	SignUp(email string, password string) (*uuid.UUID, error)
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Server struct {
	http.Handler
	render Render
	auth   Auth
}

func NewServer(auth Auth, render Render) *Server {
	server := &Server{
		auth:   auth,
		render: render,
	}

	router := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("static"))
	staticFileServer := http.FileServer(staticDir)
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFileServer))
	router.Handle("GET /favicon.ico", http.StripPrefix("/", staticFileServer))

	router.HandleFunc("GET /", server.getHomeHandler)
	router.HandleFunc("GET /account", server.getAccountHandler)

	router.HandleFunc("GET /account/sign-in", server.signInModalHandler)
	router.HandleFunc("POST /account/sign-in", server.signInHandler)

	router.HandleFunc("GET /account/sign-up", server.signUpModalHandler)
	router.HandleFunc("POST /account/sign-up", server.signUpHandler)

	server.Handler = router
	return server
}

func (s *Server) getHomeHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	signedIn := middleware.IsSignedIn(r.Context())

	if isHtmx {
		err := s.render.RenderHomeFragment(w)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderHomePage(w, signedIn)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	if isHtmx {
		err := s.render.RenderAccountFragment(w)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderAccountPage(w)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) signInModalHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	if isHtmx {
		err := s.render.RenderSignInModal(w)
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
}

func (s *Server) signInHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	user, err := s.auth.SignIn(email, password)
	if err != nil {
		err = s.render.RenderErrorMessageFragment(w, "internal server error")
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return
	}
	if user == nil {
		err = s.render.RenderErrorMessageFragment(w, "invalid sign in: check email and password")
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return
	}

	// if yes
	cookie := getSecureCookie(user.ID)
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) signUpModalHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	if isHtmx {
		err := s.render.RenderSignUpModal(w)
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
}

func (s *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	_, err := mail.ParseAddress(email)
	if err != nil || len(password) < 8 {
		err = s.render.RenderErrorMessageFragment(w, "invalid sign up: provide email and password at least 8 characters long")
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return
	}

	newID, err := s.auth.SignUp(email, password)
	if err != nil {
		err = s.render.RenderErrorMessageFragment(w, "nfi")
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return
	}

	if newID == nil {
		err = s.render.RenderErrorMessageFragment(w, "internal server error")
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return
	}
	// if yes
	cookie := getSecureCookie(*newID)
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
	w.Header().Set("Content-Type", "text/html")
}

func getSecureCookie(userID uuid.UUID) http.Cookie {
	oneDay := 24 * time.Hour
	cookie := http.Cookie{
		Name:     "userID",
		Value:    userID.String(),
		MaxAge:   int(oneDay.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
}
