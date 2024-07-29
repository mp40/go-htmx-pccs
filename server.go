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

	router.Handle("/characters", http.HandlerFunc(server.charactersHandler))

	router.Handle("/tools", http.HandlerFunc(server.toolsHandler))
	router.Handle("/tools/shooting", http.HandlerFunc(server.toolsShootingHandler))
	router.Handle("/tools/shotguns", http.HandlerFunc(server.toolsShotgunsHandler))
	router.Handle("/tools/hand-to-hand", http.HandlerFunc(server.toolsHandToHandHandler))

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

func (s *Server) charactersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderCharactersPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderCharactersFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) toolsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderToolsPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderToolsFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) toolsShootingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderToolsShootingPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderToolsShootingFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) toolsShotgunsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderToolsShotgunsPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderToolsShotgunsFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) toolsHandToHandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := render.RenderToolsHandToHandPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		err := render.RenderToolsHandToHandFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}
