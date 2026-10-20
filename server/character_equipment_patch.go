package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type patchCharacterEquipmentRender interface {
	RenderCharacterFragment(w io.Writer, character domain.EnrichedCharacterDTO) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type patchCharacterEquipmentCharacterService interface {
	UpdateCharacterEquipment(userID uuid.UUID, characterID uuid.UUID, equipmentID int, load string, qty int) (*domain.EnrichedCharacterDTO, error)
}

type patchCharacterEquipmentIdentity interface {
	GetUserID(r *http.Request) *uuid.UUID
}

func patchCharacterEquipmentHandler(render patchCharacterEquipmentRender, characterService patchCharacterEquipmentCharacterService, identity patchCharacterEquipmentIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("patch character equipment error", "err", "character id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed character id"); err != nil {
				slog.Error("patch character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}

		rawEquipmentID := r.PathValue("equipmentId")
		equipmentID, err := strconv.Atoi(rawEquipmentID)
		if err != nil {
			slog.Error("patch character equipment error", "err", "equipment id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed equipment id"); err != nil {
				slog.Error("patch character equipment error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}

		load := r.URL.Query().Get("load")
		if load != "primary" && load != "secondary" {
			slog.Error("patch character equipment error, incorrect load name", "err", load)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: primary or secondary load name required"); err != nil {
				slog.Error("patch character equipment error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("patch character equipment error, user id is nil")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := render.RenderErrorMessageFragment(w, "unauthorised: sign in to update character equipment"); err != nil {
				slog.Error("patch character equipment error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.Error("post character, patch form error", "err", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		rawQty := strings.TrimSpace(r.PostForm.Get("qty"))
		qty, err := strconv.Atoi(rawQty)
		if err != nil {
			slog.Error("patch character equipment: malformed qty", "err", rawQty)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed qty value"); err != nil {
				slog.Error("patch character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		if qty < 0 {
			slog.Error("patch character equipment: negative qty value", "err", qty)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: negative qty value"); err != nil {
				slog.Error("patch character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		c, err := characterService.UpdateCharacterEquipment(*userID, characterID, equipmentID, load, qty)
		if err != nil {
			slog.Error("patch character equipment internal server error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("patch character equipment error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		if c == nil {
			slog.Error("post character equipment error, character is nil")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("post character equipment error", "err", err)
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
