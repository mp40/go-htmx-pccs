package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type getGearPageRender interface {
	RenderGearPage(w io.Writer, signedIn bool) error
	RenderGearFragment(w io.Writer) error
}

type getGearPageIdentity interface {
	IsSignedIn(r *http.Request) bool
}

func getGearHandler(render getGearPageRender, identity getGearPageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderGearFragment(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderGearPage(w, signedIn)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
