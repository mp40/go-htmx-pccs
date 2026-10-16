package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type postCharacterEquipmentRender interface {
	RenderCharacterFragment(w io.Writer, character domain.EnrichedCharacterDTO) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type postCharacterEquipmentCharacterService interface {
	AddCharacterEquipment(userID uuid.UUID, characterID uuid.UUID, equipmentID int, load string) (*domain.EnrichedCharacterDTO, error)
}

type postCharacterEquipmentIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func postCharacterEquipmentHandler(render postCharacterEquipmentRender, characterService postCharacterEquipmentCharacterService, identity postCharacterEquipmentIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("post character equipment error", "err", "charcacter id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed character id"); err != nil {
				slog.Error("post character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			fmt.Println("WTF")
			slog.Error("post character equipment: user id nil")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := render.RenderErrorMessageFragment(w, "unauthorised: sign in to add character equipment"); err != nil {
				slog.Error("post character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("post character, parse form error", "err", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		load := strings.TrimSpace(r.PostForm.Get("load"))
		if load != "primary" && load != "secondary" {
			slog.Error("post character equipment: malformed load name")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed load name"); err != nil {
				slog.Error("post character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		rawID := strings.TrimSpace(r.PostForm.Get("equipmentId"))
		equipmentID, err := strconv.Atoi(rawID)
		if err != nil {
			slog.Error("post character equipment: malformed equidment id")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed equipment id"); err != nil {
				slog.Error("post character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		c, err := characterService.AddCharacterEquipment(*userID, characterID, equipmentID, load)
		if err != nil {
			slog.Error("post character equipment error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("post character equipment error", "err", err)
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

		err = render.RenderCharacterFragment(w, *c)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}
