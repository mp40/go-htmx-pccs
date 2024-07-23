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

	router.Handle("/", http.HandlerFunc(server.getRootHandler))

	server.Handler = router
	return server
}

func (s *Server) getRootHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderIndexPage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
