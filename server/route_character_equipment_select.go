package server

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/state"
)

type getSelectEquipmentRender interface {
	RenderEquipmentSelect(w io.Writer, characterID uuid.UUID, load string, equipment []state.Equipment) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type getSelectEquipmentIdentity interface {
	IsSignedIn(r *http.Request) bool
	GetUserID(r *http.Request) *uuid.UUID
}

type getSelectEquipmentGearService interface {
	GetEquipment() ([]state.Equipment, error)
}

func getCharacterEquipmentSelectHandler(render getSelectEquipmentRender, identity getSelectEquipmentIdentity, gearService getSelectEquipmentGearService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		load := r.URL.Query().Get("load")
		if load != "primary" && load != "secondary" {
			slog.Error("get character equipment select error, incorrect load name", "err", load)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: primary or secondary load name required"); err != nil {
				slog.Error("get character equipment selec error", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		userID := identity.GetUserID(r)
		if userID == nil {
			slog.Error("get character equipment select: user id is nil", "path", r.URL.Path)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			if err := render.RenderErrorMessageFragment(w, "unauthorised: sign in to select character equipment"); err != nil {
				slog.Error("character equipment select error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		rawCharacterID := r.PathValue("id")
		characterID, err := uuid.Parse(rawCharacterID)
		if err != nil {
			slog.Error("get character equipment select", "err", "character id is malformed")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed character id"); err != nil {
				slog.Error("character equipment select error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		equipment, err := gearService.GetEquipment()
		if err != nil {
			slog.Error("get character equipment select", "err", "get equipment error")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("character equipment select error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = render.RenderEquipmentSelect(w, characterID, load, equipment)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
