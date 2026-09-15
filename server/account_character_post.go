package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type postCharacterRender interface {
	RenderCharacter(w io.Writer, character domain.CharacterDTO) error
	RenderErrorListFragment(w io.Writer, errors []string) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type postCharacterService interface {
	AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error)
}

type postCharacterIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func postCharacterHandler(render postCharacterRender, characterService postCharacterService, identity postCharacterIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("get user id error", "err", "user id is nil")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if err := r.ParseForm(); err != nil {
			slog.Error("post character, parse form error", "err", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		c, problems := domain.ParseRawCharacter(r.PostForm)
		if len(problems) > 0 {
			slog.Warn("post character, invalid data submitted", "problems", problems)
			errorList := []string{}
			for k, v := range problems {
				errorList = append(errorList, fmt.Sprintf("%s %s", k, v))
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			err := render.RenderErrorListFragment(w, errorList)
			if err != nil {
				slog.Error("parse character error list", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		character, err := characterService.AddCharacter(*userID, c)
		if err != nil {
			slog.Error("add character error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("render error message", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		// once response body is written to can't set header - call before otherwise 200
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
