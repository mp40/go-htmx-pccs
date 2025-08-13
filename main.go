package main

import (
	"log"
	"net/http"

	"github.com/mp40/go-htmx-pccs/render"
)

func main() {
	r := render.NewRenderService()
	server := NewServer(r)
	log.Fatal(http.ListenAndServe(":5050", server))
}
