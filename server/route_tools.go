package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type getToolsPageRender interface {
	RenderToolsPage(w io.Writer, signedIn bool) error
	RenderToolsFragment(w io.Writer) error
}

type getToolsPageIdentity interface {
	IsSignedIn(r *http.Request) bool
}

func getToolsHandler(render getToolsPageRender, identity getToolsPageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		if isHtmx {
			err := render.RenderToolsFragment(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderToolsPage(w, signedIn)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/html")
	})
}
