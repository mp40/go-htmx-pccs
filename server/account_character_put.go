package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type putCharacterRender interface {
	RenderAccountFragment(w io.Writer, characters []domain.CharacterDTO) error
	RenderErrorListFragment(w io.Writer, errors []string) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type putCharacterService interface {
	EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error)
	GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error)
}

type putCharacterIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func putCharacterHandler(render putCharacterRender, characterService putCharacterService, identity putCharacterIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("edit charcater error", "err", "charcacter id is malformed")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("edit character error", "err", "user id is nil")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("put character, parse form error", "err", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		c, problems := parseRawCharacter(r.PostForm)
		if len(problems) > 0 {
			slog.Warn("put character, invalid data submitted", "problems", problems)
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

		// TODO character nil check
		_, err = characterService.EditCharacter(*userID, characterID, c)
		if err != nil {
			slog.Error("edit character error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("render error message", "err", err)
				return
			}
			return
		}

		characters, err := characterService.GetCharactersByUserID(*userID)
		if err != nil {
			slog.Error("get characters error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("render error message", "err", err)
				return
			}
			return
		}

		w.Header().Set("HX-Push-Url", "/account")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = render.RenderAccountFragment(w, characters)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
