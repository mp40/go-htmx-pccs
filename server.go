package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type Render interface {
	RenderHomePage(w io.Writer, signedIn bool) error
	RenderHomeFragment(w io.Writer) error
	RenderAccountPage(w io.Writer, signedIn bool) error
	RenderAccountFragment(w io.Writer, signedIn bool) error
	RenderReferencePage(w io.Writer) error
	RenderReferenceFragment(w io.Writer) error
	RenderToolsPage(w io.Writer) error
	RenderToolsFragment(w io.Writer) error
	RenderSignInModal(w io.Writer) error
	RenderSignUpModal(w io.Writer) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
	RenderCharacter(w io.Writer, character store.Character) error
}

type Auth interface {
	SignIn(email string, password string) (*store.User, error)
	SignUp(email string, password string) (*uuid.UUID, error)
}

type Session interface {
	AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error)
	DeleteSessionByID(ID uuid.UUID) error
}

type CharacterService interface {
	AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*store.Character, error)
}

type Identity interface {
	GetUserID(r *http.Request) *uuid.UUID
	IsSignedIn(r *http.Request) bool
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Server struct {
	http.Handler
	render           Render
	auth             Auth
	session          Session
	characterService CharacterService
	identity         Identity
}

func NewServer(auth Auth, session Session, identity Identity, characterService CharacterService, render Render) *Server {
	server := &Server{
		auth:             auth,
		session:          session,
		identity:         identity,
		render:           render,
		characterService: characterService,
	}

	router := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("static"))
	staticFileServer := http.FileServer(staticDir)
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFileServer))
	router.Handle("GET /favicon.ico", http.StripPrefix("/", staticFileServer))

	router.HandleFunc("GET /", server.getHomeHandler)
	router.HandleFunc("GET /account", server.getAccountHandler)

	router.Handle("GET /reference", getReferenceHandler(server.render))
	router.Handle("GET /tools", getToolsHandler(server.render))

	router.HandleFunc("GET /account/sign-in", server.signInModalHandler)
	router.HandleFunc("POST /account/sign-in", server.signInHandler)

	router.HandleFunc("GET /account/sign-up", server.signUpModalHandler)
	router.HandleFunc("POST /account/sign-up", server.signUpHandler)

	router.HandleFunc("DELETE /account/sign-out", server.signOutHandler)

	router.Handle("POST /account/characters", handlePostCharacter(server.render, server.characterService, server.identity))

	server.Handler = router
	return server
}

func (s *Server) getHomeHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	signedIn := s.identity.IsSignedIn(r)

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
	signedIn := s.identity.IsSignedIn(r)

	if isHtmx {
		err := s.render.RenderAccountFragment(w, signedIn)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderAccountPage(w, signedIn)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func getReferenceHandler(render Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)

		if isHtmx {
			err := render.RenderReferenceFragment(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderReferencePage(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
	})
}

func getToolsHandler(render Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)

		if isHtmx {
			err := render.RenderToolsFragment(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderToolsPage(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
	})
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
	now := time.Now()
	sessionID, err := s.session.AddSession(user.ID, now.Add(12*time.Hour))
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
	cookie := getSecureCookie(sessionID)
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
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
	now := time.Now()
	sessionID, err := s.session.AddSession(*newID, now.Add(12*time.Hour))
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
	cookie := getSecureCookie(sessionID)
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) signOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	sessionID, err := uuid.Parse(cookie.Value)
	if err != nil {
		redirectWithExpiredCookie(w)
		return
	}

	err = s.session.DeleteSessionByID(sessionID)
	if err != nil {
		slog.Error("delete session error", "err", err)
	}

	redirectWithExpiredCookie(w)
}

func handlePostCharacter(render Render, characterService CharacterService, identity Identity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("get user id error", "err", "user id is nil")
			// handle nil user id somehow
			return
		}

		if err := r.ParseForm(); err != nil {
			slog.Error("post character, parse form error", "err", err)
			// do something
			return
		}

		c, problems := parseRawCharacter(r.PostForm)
		if len(problems) > 0 {
			slog.Warn("post character, invalid data submitted", "problems", problems)
			// do something in UI
			return
		}

		character, err := characterService.AddCharacter(*userID, c)
		if err != nil {
			slog.Error("post character, add character error", "err", err)
			// handle err some how
			return
		}

		// return render character fragment
		err = render.RenderCharacter(w, *character)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusCreated)
	})
}

func redirectWithExpiredCookie(w http.ResponseWriter) {
	expiredCookie := removeSecureCookie()
	http.SetCookie(w, &expiredCookie)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func getSecureCookie(sessionID uuid.UUID) http.Cookie {
	oneDay := 24 * time.Hour
	cookie := http.Cookie{
		Name:     "sessionID",
		Value:    sessionID.String(),
		MaxAge:   int(oneDay.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	return cookie
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

func parseRawCharacter(form url.Values) (domain.RawCharacter, map[string]string) {
	problems := map[string]string{}

	name := strings.TrimSpace(form.Get("name"))
	rawStr := strings.TrimSpace(form.Get("str"))
	rawItel := strings.TrimSpace(form.Get("int"))
	rawWil := strings.TrimSpace(form.Get("wil"))
	rawHlt := strings.TrimSpace(form.Get("hlt"))
	rawAgi := strings.TrimSpace(form.Get("agi"))
	rawTch := strings.TrimSpace(form.Get("tch"))
	rawGunCombatLevel := strings.TrimSpace(form.Get("gun_combat_level"))
	rawHandToHandLevel := strings.TrimSpace(form.Get("hand_to_hand_level"))

	str, err := strconv.Atoi(rawStr)
	if err != nil {
		problems["str"] = "is not a number"
	} else if str < 1 || str > 21 {
		problems["str"] = "is invalid (must be between 1 and 21)"
	}

	intel, err := strconv.Atoi(rawItel)
	if err != nil {
		problems["int"] = "is not a number"
	} else if intel < 1 || intel > 21 {
		problems["int"] = "is invalid (must be between 1 and 21)"
	}

	wil, err := strconv.Atoi(rawWil)
	if err != nil {
		problems["wil"] = "is not a number"
	} else if wil < 1 || wil > 21 {
		problems["wil"] = "is invalid (must be between 1 and 21)"
	}

	hlt, err := strconv.Atoi(rawHlt)
	if err != nil {
		problems["hlt"] = "is not a number"
	} else if hlt < 1 || hlt > 21 {
		problems["hlt"] = "is invalid (must be between 1 and 21)"
	}

	agi, err := strconv.Atoi(rawAgi)
	if err != nil {
		problems["agi"] = "is not a number"
	} else if agi < 1 || agi > 21 {
		problems["agi"] = "is invalid (must be between 1 and 21)"
	}

	tch, err := strconv.Atoi(rawTch)
	if err != nil {
		problems["tch"] = "is not a number"
	} else if tch < 1 || tch > 21 {
		problems["tch"] = "is invalid (must be between 1 and 21)"
	}

	gunCombatLevel, err := strconv.Atoi(rawGunCombatLevel)
	if err != nil {
		problems["gun_combat_level"] = "is not a number"
	} else if gunCombatLevel < 0 || gunCombatLevel > 20 {
		problems["gun_combat_level"] = "is invalid (must be between 0 and 20)"
	}

	handToHandLevel, err := strconv.Atoi(rawHandToHandLevel)
	if err != nil {
		problems["hand_to_hand_level"] = "is not a number"
	} else if handToHandLevel < 0 || handToHandLevel > 20 {
		problems["hand_to_hand_level"] = "is invalid (must be between 0 and 20)"
	}

	c := domain.RawCharacter{
		Name:            name,
		Str:             str,
		Int:             intel,
		Wil:             wil,
		Hlt:             hlt,
		Agi:             agi,
		Tch:             tch,
		GunCombatLevel:  gunCombatLevel,
		HandToHandLevel: handToHandLevel,
	}

	return c, problems
}
