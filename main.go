package main

import (
	"log"
	"net/http"

	"github.com/mp40/go-htmx-pccs/auth"
	"github.com/mp40/go-htmx-pccs/middleware"
	"github.com/mp40/go-htmx-pccs/render"
)

func main() {

	a := auth.NewAuthService()
	r := render.NewRenderService()

	server := NewServer(a, r)

	m := middleware.NewMiddlewareService()
	serverWithMiddleware := m.AuthMiddleware(server)

	log.Fatal(http.ListenAndServe(":5050", serverWithMiddleware))
}
