package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/state"
	"github.com/mp40/go-htmx-pccs/store"
)

type Auth interface {
	SignIn(email string, password string) (*store.User, error)
	SignUp(email string, password string) (*uuid.UUID, error)
}

type Session interface {
	AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error)
	DeleteSessionByID(ID uuid.UUID) error
}

type CharacterService interface {
	AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error)
	EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacterEdit domain.RawCharacter) (*domain.CharacterDTO, error)
	GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error)
	GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error)
	DeleteUserCharacterByID(userID uuid.UUID, charcaterID uuid.UUID) error
}

type Identity interface {
	GetUserID(r *http.Request) *uuid.UUID
	IsSignedIn(r *http.Request) bool
}

type GearService interface {
	GetEquipment() ([]state.Equipment, error)
}

func NewServer(auth Auth, session Session, identity Identity, characterService CharacterService, gearService GearService, render *render.Render) http.Handler {
	router := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("static"))
	staticFileServer := http.FileServer(staticDir)
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFileServer))
	router.Handle("GET /favicon.ico", http.StripPrefix("/", staticFileServer))

	router.Handle("GET /", getHomeHandler(render, identity))
	router.Handle("GET /account", getAccountHandler(render, identity, characterService))

	router.Handle("GET /reference", getReferenceHandler(render, identity))

	router.Handle("GET /tools", getToolsHandler(render, identity))
	router.Handle("GET /tools/shotguns/spread", getToolsShotgunSpreadHandler(render))

	router.Handle("GET /gear", getGearHandler(render, identity, gearService))

	router.Handle("GET /account/sign-in", getSignInHandler(render))
	router.Handle("POST /account/sign-in", postSignInHandler(render, auth, session))

	router.Handle("GET /account/sign-up", getSignUpHandler(render))
	router.Handle("POST /account/sign-up", postSignUpHandler(render, auth, session))

	router.Handle("DELETE /account/sign-out", deleteSignOutHandler(session))

	router.Handle("POST /account/characters", postCharacterHandler(render, characterService, identity))
	router.Handle("PUT /account/characters/{id}", putCharacterHandler(render, characterService, identity))
	router.Handle("DELETE /account/characters/{id}", deleteCharacterHandler(render, characterService, identity))
	router.Handle("GET /account/characters/{id}/edit", getCharacterEditHandler(render, characterService, identity))

	router.Handle("GET /character/{id}", getCharacterHandler(render, characterService, identity))

	return router
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
