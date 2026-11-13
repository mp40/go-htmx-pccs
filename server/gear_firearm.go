package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mp40/go-htmx-pccs/domain"
)

type getGearFirearmRender interface {
	RenderEnrichedFirearm(w io.Writer, firearm domain.EnrichedFirearmDTO) error
	RenderErrorMessageFragment(w io.Writer, msg string) error
}

type getGearFirearmGearService interface {
	GetEnrichedFirearmByID(id int) (*domain.EnrichedFirearmDTO, error)
}

func getGearFirearmHandler(render getGearFirearmRender, gearService getGearFirearmGearService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawFirearmID := r.PathValue("id")
		firearmID, err := strconv.Atoi(rawFirearmID)
		if err != nil {
			slog.Error("get gear firearm error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			if err := render.RenderErrorMessageFragment(w, "bad request: malformed firearm id"); err != nil {
				slog.Error("get gear firearm error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		firearm, err := gearService.GetEnrichedFirearmByID(firearmID)
		if err != nil {
			slog.Error("get gear firearm error", "err", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			if err := render.RenderErrorMessageFragment(w, "internal server error"); err != nil {
				slog.Error("get gear firearm error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		if firearm == nil {
			slog.Warn("get gear firearm error, firearm is nil", "firearmId", firearmID)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			if err := render.RenderErrorMessageFragment(w, "not found: firearm data not found"); err != nil {
				slog.Error("get gear firearm error", "err", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = render.RenderEnrichedFirearm(w, *firearm)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	})
}
