package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type getCharacterPageRender interface {
	RenderCharacterPage(w io.Writer, signedIn bool, character domain.CharacterDTO) error
	RenderCharacterFragment(w io.Writer, character domain.CharacterDTO) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type getCharacterPageIdentity interface {
	IsSignedIn(r *http.Request) bool
	GetUserID(r *http.Request) *uuid.UUID
}

type getCharacterPageCharacters interface {
	GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error)
}

func getCharacterHandler(render getCharacterPageRender, characters getCharacterPageCharacters, identity getCharacterPageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("get character page: user id is nil", "path", r.URL.Path)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := render.RenderErrorMessageFragment(w, "unauthorised: sign in to edit character"); err != nil {
				slog.Error("character page error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("get character page", "err", "character id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed character id"); err != nil {
				slog.Error("character page error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		c, err := characters.GetUserCharacterByID(*userID, characterID)
		if err != nil {
			slog.Error("get character page", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("character page error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}
		if c == nil {
			slog.Warn("get character page: character not found", "character_id", characterID, "user_id", *userID)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			if err := render.RenderErrorMessageFragment(w, "character not found"); err != nil {
				slog.Error("character page error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err = render.RenderCharacterFragment(w, *c)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err = render.RenderCharacterPage(w, signedIn, *c)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
