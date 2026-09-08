package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type deleteCharacterRender interface {
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type deleteCharacterService interface {
	DeleteUserCharacterByID(userID uuid.UUID, charcaterID uuid.UUID) error
}

type deleteCharacterIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func deleteCharacterHandler(render deleteCharacterRender, characterService deleteCharacterService, identity deleteCharacterIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("delete character error", "err", "character id is malformed")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("delete character error", "err", "user id is nil")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		err = characterService.DeleteUserCharacterByID(*userID, characterID)
		if err != nil {
			slog.Error("delete character error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			err = render.RenderErrorMessageFragment(w, "internal server error")
			if err != nil {
				slog.Error("render error message", "err", err)
				return
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	})
}
