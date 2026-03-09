package main

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Render interface {
	RenderHomePage(w io.Writer) error
	RenderHomeFragment(w io.Writer) error
	RenderAccountPage(w io.Writer) error
	RenderAccountFragment(w io.Writer) error
	RenderAccountSignInFailureFragment(w io.Writer) error
	RenderSignInModal(w io.Writer) error
	RenderAccountSignUpFailureFragment(w io.Writer, msg string) error
	RenderSignUpModal(w io.Writer) error
}

type Auth interface {
	SignIn(email string, password string) error
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
	if isHtmx {
		err := s.render.RenderHomeFragment(w)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderHomePage(w)
		if err != nil {
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderAccountPage(w)
		if err != nil {
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// ???
		// what to do if not htmx?
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) signInHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := s.auth.SignIn(email, password)

	// if no - something
	if err != nil {
		err = s.render.RenderAccountSignInFailureFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// if yes
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
	}
	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) signUpModalHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	if isHtmx {
		err := s.render.RenderSignUpModal(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// ???
		// what to do if not htmx?
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) signUpHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	newID, err := s.auth.SignUp(email, password)
	if err != nil {
		err = s.render.RenderAccountSignUpFailureFragment(w, "nfi")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		return
	}

	if newID == nil {
		err = s.render.RenderAccountSignUpFailureFragment(w, "nfi")
		if err != nil {
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
