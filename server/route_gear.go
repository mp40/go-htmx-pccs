package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mp40/go-htmx-pccs/state"
)

type getGearPageRender interface {
	RenderGearPage(w io.Writer, signedIn bool, equipment []state.Equipment) error
	RenderGearFragment(w io.Writer, equipment []state.Equipment) error
}

type getGearPageIdentity interface {
	IsSignedIn(r *http.Request) bool
}

type getGearPageGearService interface {
	GetEquipment() ([]state.Equipment, error)
}

func getGearHandler(render getGearPageRender, identity getGearPageIdentity, gearService getGearPageGearService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		equipment, err := gearService.GetEquipment()
		if err != nil {
			slog.Error("state get equipment", "err", err)
			// TODO decide what to do on equip error
			// for now empty list is OK
			// in future maybe also error msg for UI
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderGearFragment(w, equipment)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderGearPage(w, signedIn, equipment)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
