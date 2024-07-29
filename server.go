package main

import (
	"net/http"
	"path/filepath"

	"github.com/mp40/go-htmx-pccs/render"
)

type Server struct {
	http.Handler
}

func NewServer() *Server {
	server := new(Server)

	router := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("static"))
	staticFileServer := http.FileServer(staticDir)
	router.Handle("/static/", http.StripPrefix("/static/", staticFileServer))
	router.Handle("/favicon.ico", http.StripPrefix("/", staticFileServer))

	router.Handle("/", http.HandlerFunc(server.homeHandler))
	router.Handle("/shooting", http.HandlerFunc(server.shootingHandler))

	server.Handler = router
	return server
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderHomePage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderHomeFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) shootingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderShootingPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderShootingFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}
