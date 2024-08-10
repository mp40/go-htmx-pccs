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
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFileServer))
	router.Handle("GET /favicon.ico", http.StripPrefix("/", staticFileServer))

	router.HandleFunc("GET /", http.HandlerFunc(server.getHomeHandler))
	router.HandleFunc("POST /", http.HandlerFunc(server.postHomeHandler))

	router.HandleFunc("GET /characters", http.HandlerFunc(server.getCharactersHandler))
	router.HandleFunc("POST /characters", http.HandlerFunc(server.postCharactersHandler))

	router.HandleFunc("GET /tools", http.HandlerFunc(server.getToolsHandler))
	router.HandleFunc("POST /tools", http.HandlerFunc(server.postToolsHandler))
	router.HandleFunc("GET /tools/shooting", http.HandlerFunc(server.getToolsShootingHandler))
	router.HandleFunc("POST /tools/shooting", http.HandlerFunc(server.postToolsShootingHandler))
	router.HandleFunc("GET /tools/shooting/calculator", http.HandlerFunc(server.getToolsShootingCalculatorHandler))
	router.HandleFunc("GET /tools/shotguns", http.HandlerFunc(server.getToolsShotgunsHandler))
	router.HandleFunc("POST /tools/shotguns", http.HandlerFunc(server.postToolsShotgunsHandler))
	router.HandleFunc("GET /tools/hand-to-hand", http.HandlerFunc(server.getToolsHandToHandHandler))
	router.HandleFunc("POST /tools/hand-to-hand", http.HandlerFunc(server.postToolsHandToHandHandler))

	server.Handler = router
	return server
}

func (s *Server) getHomeHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderHomePage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) postHomeHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderHomeFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getCharactersHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderCharactersPage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) postCharactersHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderCharactersFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getToolsHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsPage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) postToolsHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getToolsShootingHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsShootingPage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) postToolsShootingHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsShootingFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getToolsShootingCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsShootingCalculatorFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getToolsShotgunsHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsShotgunsPage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) postToolsShotgunsHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsShotgunsFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getToolsHandToHandHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsHandToHandPage(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) postToolsHandToHandHandler(w http.ResponseWriter, r *http.Request) {
	err := render.RenderToolsHandToHandFragment(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
}
