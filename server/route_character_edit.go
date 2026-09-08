package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type getCharacterEditPageRender interface {
	RenderCharacterEditPage(w io.Writer, signedIn bool, character domain.CharacterDTO) error
	RenderCharacterEditFragment(w io.Writer, character domain.CharacterDTO) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type getCharacterEditPageIdentity interface {
	IsSignedIn(r *http.Request) bool
	GetUserID(r *http.Request) *uuid.UUID
}

type getCharacterEditPageCharacters interface {
	GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error)
}

func getCharacterEditHandler(render getCharacterEditPageRender, characters getCharacterEditPageCharacters, identity getCharacterEditPageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("get edit character page: user id is nil", "path", r.URL.Path)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := render.RenderErrorMessageFragment(w, "unauthorised: sign in to edit character"); err != nil {
				slog.Error("edit character error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("get edit character page", "err", "character id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed character id"); err != nil {
				slog.Error("edit character error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		c, err := characters.GetUserCharacterByID(*userID, characterID)
		if err != nil {
			slog.Error("get edit character page", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("edit character error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}
		if c == nil {
			slog.Warn("get edit character page: character not found", "character_id", characterID, "user_id", *userID)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			if err := render.RenderErrorMessageFragment(w, "character not found"); err != nil {
				slog.Error("edit character error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderCharacterEditFragment(w, *c)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err = render.RenderCharacterEditPage(w, signedIn, *c)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	})
}
