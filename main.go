package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mp40/go-htmx-pccs/auth"
	"github.com/mp40/go-htmx-pccs/middleware"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/store"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", "./pccs_store.db")
	if err != nil {
		log.Fatalf("init store db error: %v", err)
	}
	defer db.Close()
	// fine for now but need to do this only if in local dev mode once deployed
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY NOT NULL, email TEXT UNIQUE NOT NULL, hash TEXT NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL)")
	if err != nil {
		log.Fatalf("db create table error: %v", err)
	}

	s := store.NewStoreService(db)
	a := auth.NewAuthService(s)
	r := render.NewRenderService()

	server := NewServer(a, r)

	m := middleware.NewMiddlewareService(s)
	serverWithMiddleware := m.AuthMiddleware(server)

	slog.Info("server running", "on", "http://localhost:5050")
	log.Fatal(http.ListenAndServe(":5050", serverWithMiddleware))
}
