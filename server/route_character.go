package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type getCharacterPageRender interface {
	RenderCharacterPage(w io.Writer, signedIn bool, character domain.CharacterDTO) error
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
		signedIn := identity.IsSignedIn(r)

		userID := identity.GetUserID(r)
		if userID == nil {
			// TODO finish once failing test to pass
			return
		}

		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			// TODO finsih once failing test to pass
			return
		}

		c, err := characters.GetUserCharacterByID(*userID, characterID)
		if err != nil {
			// TODO finish once failing test to pass
			return
		}
		if c == nil {
			// TODO finsih once failing test to pass
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = render.RenderCharacterPage(w, signedIn, *c)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
