package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type getReferencePageRender interface {
	RenderReferencePage(w io.Writer, signedIn bool) error
	RenderReferenceFragment(w io.Writer) error
}

type getReferencePageIdentity interface {
	IsSignedIn(r *http.Request) bool
}

func getReferenceHandler(render getReferencePageRender, identity getReferencePageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderReferenceFragment(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderReferencePage(w, signedIn)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
