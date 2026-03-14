package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mp40/go-htmx-pccs/auth"
	"github.com/mp40/go-htmx-pccs/data"
	"github.com/mp40/go-htmx-pccs/middleware"
	"github.com/mp40/go-htmx-pccs/render"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./data/app.db")
	if err != nil {
		log.Fatalf("init db error: %v", err)
	}

	d := data.NewDataService(db)
	a := auth.NewAuthService(d)
	r := render.NewRenderService()

	server := NewServer(a, r)

	m := middleware.NewMiddlewareService(d)
	serverWithMiddleware := m.AuthMiddleware(server)

	log.Fatal(http.ListenAndServe(":5050", serverWithMiddleware))
	slog.Info("server running on :5050")
}
