package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type getCharacterEditPageRender interface {
	RenderCharacterEditPage(w io.Writer, signedIn bool, character domain.CharacterDTO) error
}

type getCharacterEditPageIdentity interface {
	IsSignedIn(r *http.Request) bool
	GetUserID(r *http.Request) *uuid.UUID
}

type getCharacterEditPageCharacters interface {
	GetCharacterByID(ID uuid.UUID) (*domain.CharacterDTO, error)
}

func getCharacterEditHandler(render getCharacterEditPageRender, characters getCharacterEditPageCharacters, identity getCharacterEditPageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO HTMX check full page vs fragment
		signedIn := identity.IsSignedIn(r)
		// TODO handle not signed in

		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("edit charcater error", "err", "charcacter id is malformed")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		c, err := characters.GetCharacterByID(characterID)
		if err != nil {
			slog.Error("get character by id", "err", err)
			// TODOD handle error
			return
		}
		if c == nil {
			// warn or err?
			slog.Warn("character not found by id", "id", characterID)
			// TODO handle nil character
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = render.RenderCharacterEditPage(w, signedIn, *c)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
