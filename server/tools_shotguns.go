package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mp40/go-htmx-pccs/domain/pccs"
)

type getToolsShotgunSpreadHandlerRender interface {
	RenderShotgunPelletHitsFragment(w io.Writer, hits []int) error
	RenderErrorListFragment(w io.Writer, errors []string) error
}

func getToolsShotgunSpreadHandler(render getToolsShotgunSpreadHandlerRender) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawLocation := r.URL.Query().Get("initialLocation")
		rawSalm := r.URL.Query().Get("salm")
		rawHits := r.URL.Query().Get("hitCount")

		errors := []string{}

		initialLocation, err := strconv.Atoi(rawLocation)
		if err != nil {
			errors = append(errors, "initial hit location value malformed")
		}
		salm, err := strconv.Atoi(rawSalm)
		if err != nil {
			errors = append(errors, "SALM value malformed")
		}
		additionalHits, err := strconv.Atoi(rawHits)
		if err != nil {
			errors = append(errors, "additional hits value malformed")
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if len(errors) > 0 {
			err = render.RenderErrorListFragment(w, errors)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		hits := pccs.GenerateRandomHitsWithinSpread(initialLocation, salm, additionalHits, pccs.RollRangeAdvancedOpen)

		err = render.RenderShotgunPelletHitsFragment(w, hits)
		if err != nil {
			slog.Error("server error rendering", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
