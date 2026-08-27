package server

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
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
	EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacterEdit domain.RawCharacterEdit) (*domain.CharacterDTO, error)
	GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error)
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

func parseRawCharacterEdit(form url.Values) (domain.RawCharacterEdit, map[string]string) {
	problems := map[string]string{}

	name := strings.TrimSpace(form.Get("name"))
	rawStr := strings.TrimSpace(form.Get("str"))
	rawItel := strings.TrimSpace(form.Get("int"))
	rawWil := strings.TrimSpace(form.Get("wil"))
	rawHlt := strings.TrimSpace(form.Get("hlt"))
	rawAgi := strings.TrimSpace(form.Get("agi"))
	rawTch := strings.TrimSpace(form.Get("tch"))
	rawGunCombatLearningPoints := strings.TrimSpace(form.Get("gun_combat_learning_points"))
	rawHandToHandLearningPoints := strings.TrimSpace(form.Get("hand_to_hand_learning_points"))

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

	gunCombatLearningPoints, err := strconv.ParseFloat(rawGunCombatLearningPoints, 32)
	if err != nil {
		problems["gun_combat_learning_points"] = "is not a number"
	} else if gunCombatLearningPoints < 0 || gunCombatLearningPoints > 1834 {
		problems["gun_combat_learning_points"] = "is invalid (must be between 0 and 1834)"
	}

	handToHandLearningPoints, err := strconv.ParseFloat(rawHandToHandLearningPoints, 32)
	if err != nil {
		problems["hand_to_hand_learning_points"] = "is not a number"
	} else if handToHandLearningPoints < 0 || handToHandLearningPoints > 1834 {
		problems["hand_to_hand_learning_points"] = "is invalid (must be between 0 and 1834)"
	}

	c := domain.RawCharacterEdit{
		Name:                     name,
		Str:                      str,
		Int:                      intel,
		Wil:                      wil,
		Hlt:                      hlt,
		Agi:                      agi,
		Tch:                      tch,
		GunCombatLearningPoints:  float32(gunCombatLearningPoints),
		HandToHandLearningPoints: float32(handToHandLearningPoints),
	}

	return c, problems
}
