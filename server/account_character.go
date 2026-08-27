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

		c, problems := parseRawCharacter(r.PostForm)
		if len(problems) > 0 {
			slog.Warn("post character, invalid data submitted", "problems", problems)
			errorList := []string{}
			for k, v := range problems {
				errorList = append(errorList, fmt.Sprintf("%s %s", k, v))
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorListFragment(w, errorList); err != nil {
				slog.Error("parse character error list", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		character, err := characterService.AddCharacter(*userID, c)
		if err != nil {
			slog.Error("add character error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			err = render.RenderErrorListFragment(w, []string{"internal server error"})
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

type putCharacterRender interface {
	RenderCharacter(w io.Writer, character domain.CharacterDTO) error
	RenderErrorListFragment(w io.Writer, errors []string) error
}

type putCharacterService interface {
	EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacter domain.RawCharacterEdit) (*domain.CharacterDTO, error)
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

		c, problems := parseRawCharacterEdit(r.PostForm)
		if len(problems) > 0 {
			slog.Warn("put character, invalid data submitted", "problems", problems)
			errorList := []string{}
			for k, v := range problems {
				errorList = append(errorList, fmt.Sprintf("%s %s", k, v))
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorListFragment(w, errorList); err != nil {
				slog.Error("parse character error list", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		character, err := characterService.EditCharacter(*userID, characterID, c)
		if err != nil {
			slog.Error("edit character error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			err = render.RenderErrorListFragment(w, []string{"internal server error"})
			if err != nil {
				slog.Error("render error message", "err", err)
				return
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// TODO return render character fragment
		// change placeholder? as that one is just for one line summary of character
		err = render.RenderCharacter(w, *character)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
