package server

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

type renderSignInModal interface {
	RenderSignInModal(w io.Writer) error
}

func getSignInHandler(render renderSignInModal) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmxHeader := r.Header.Get("HX-Request")
		isHtmx, _ := strconv.ParseBool(htmxHeader)
		if isHtmx {
			err := render.RenderSignInModal(w)
			if err != nil {
				slog.Error("server error rendering", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// log err
			// ???
			// what to do if not htmx?
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
	})
}
