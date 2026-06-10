package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type postCharacterRender interface {
	RenderCharacter(w io.Writer, character store.Character) error
}

type postCharacterService interface {
	AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*store.Character, error)
}

type postCharacterIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func postCharacterHandler(render postCharacterRender, characterService postCharacterService, identity postCharacterIdentity) http.Handler {
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

		// once response body is written to can't set header - call before otherwise 200
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusCreated)

		// return render character fragment
		err = render.RenderCharacter(w, *character)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
