package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type getAccountPageRender interface {
	RenderAccountPage(w io.Writer, signedIn bool) error
	RenderAccountFragment(w io.Writer, signedIn bool) error
}

type getAccountPageIdentity interface {
	IsSignedIn(r *http.Request) bool
}

func getAccountHandler(render getAccountPageRender, identity getAccountPageIdentity) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		signedIn := identity.IsSignedIn(r)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if isHtmx {
			err := render.RenderAccountFragment(w, signedIn)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			err := render.RenderAccountPage(w, signedIn)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
