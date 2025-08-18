package main

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

type Render interface {
	RenderHomePage(w io.Writer) error
	RenderHomeFragment(w io.Writer) error
	RenderAccountPage(w io.Writer) error
	RenderAccountFragment(w io.Writer) error
}

type Server struct {
	http.Handler
	render Render
}

func NewServer(render Render) *Server {
	server := &Server{
		render: render,
	}

	router := http.NewServeMux()

	staticDir := http.Dir(filepath.Join("static"))
	staticFileServer := http.FileServer(staticDir)
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFileServer))
	router.Handle("GET /favicon.ico", http.StripPrefix("/", staticFileServer))

	router.HandleFunc("GET /", server.getHomeHandler)

	server.Handler = router
	return server
}

func (s *Server) getHomeHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	if isHtmx {
		err := s.render.RenderHomeFragment(w)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderHomePage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
}

func (s *Server) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	htmxHeader := r.Header.Get("HX-Request")
	isHtmx, _ := strconv.ParseBool(htmxHeader)
	if isHtmx {
		err := s.render.RenderAccountFragment(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := s.render.RenderAccountPage(w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// func (s *Server) postHomeHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderHomeFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) getCharactersHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderCharactersPage(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) postCharactersHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderCharactersFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) getToolsHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsPage(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) postToolsHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) getToolsShootingHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsShootingPage(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) postToolsShootingHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsShootingFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) getToolsShootingCalculatorHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsShootingCalculatorFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) getToolsShotgunsHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsShotgunsPage(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) postToolsShotgunsHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsShotgunsFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) getToolsHandToHandHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsHandToHandPage(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }

// func (s *Server) postToolsHandToHandHandler(w http.ResponseWriter, r *http.Request) {
// 	err := render.RenderToolsHandToHandFragment(w)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// }
