package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type deleteCharacterEquipmentRender interface {
	RenderCharacterFragment(w io.Writer, character domain.EnrichedCharacterDTO) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type deleteCharacterEquipmentCharacterService interface {
	DeleteCharacterEquipment(userID uuid.UUID, characterID uuid.UUID, equipmentID int, load string) (*domain.EnrichedCharacterDTO, error)
}

type deleteCharacterEquipmentIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func deleteCharacterEquipmentHandler(render deleteCharacterEquipmentRender, characterService deleteCharacterEquipmentCharacterService, identity deleteCharacterEquipmentIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("delete character equipment error", "err", "character id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed character id"); err != nil {
				slog.Error("delete character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		rawEquipmentID := r.PathValue("equipmentId")
		equipmentID, err := strconv.Atoi(rawEquipmentID)
		if err != nil {
			slog.Error("delete character equipment error", "err", "equipment id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed equipment id"); err != nil {
				slog.Error("delete character equipment error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		load := r.URL.Query().Get("load")
		if load != "primary" && load != "secondary" {
			slog.Error("delete character equipment error, incorrect load name", "err", load)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: primary or secondary load name required"); err != nil {
				slog.Error("delete character equipment error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("delete character equipment error, user id is nil")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := render.RenderErrorMessageFragment(w, "unauthorised: sign in to delete character equipment"); err != nil {
				slog.Error("delete character equipment error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		c, err := characterService.DeleteCharacterEquipment(*userID, characterID, equipmentID, load)
		if err != nil {
			slog.Error("delete character equipment internal server error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("delete character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		if c == nil {
			slog.Warn("delete character equipment error, character is nil", "userId", userID, "characterId", characterID)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			if err := render.RenderErrorMessageFragment(w, "not found: character not found"); err != nil {
				slog.Error("delete character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = render.RenderCharacterFragment(w, *c)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
