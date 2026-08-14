package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type getHomePageRender interface {
	RenderHomePage(w io.Writer, signedIn bool) error
	RenderHomeFragment(w io.Writer) error
}

type getHomePageIdentity interface {
	IsSignedIn(r *http.Request) bool
}

func getHomeHandler(render getHomePageRender, identity getHomePageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderHomeFragment(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderHomePage(w, signedIn)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
