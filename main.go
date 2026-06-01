package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/mp40/go-htmx-pccs/auth"
	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/middleware"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/service"
	"github.com/mp40/go-htmx-pccs/session"
	"github.com/mp40/go-htmx-pccs/store"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	storeDB, err := sql.Open("sqlite", "./pccs_store.db")
	if err != nil {
		log.Fatalf("init store db error: %v", err)
	}
	defer storeDB.Close()
	// fine for now but need to do this only if in local dev mode once deployed
	_, err = storeDB.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY NOT NULL, email TEXT UNIQUE NOT NULL, hash TEXT NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL)")
	if err != nil {
		log.Fatalf("db create user table error: %v", err)
	}

	sessionDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("init session db error: %v", err)
	}
	defer sessionDB.Close()
	// fine for now but need to do this only if in local dev mode once deployed
	_, err = sessionDB.Exec("CREATE TABLE IF NOT EXISTS sessions (id TEXT PRIMARY KEY NOT NULL, user_id TEXT NOT NULL, created_at INTEGER NOT NULL, expires_at INTEGER NOT NULL)")
	if err != nil {
		log.Fatalf("db create session table error: %v", err)
	}

	// third SQLite db to be called state
	// a read only db ships with app for pccs game data

	storeService := store.NewStoreService(storeDB)
	sessionService := session.NewSessionService(sessionDB)
	authService := auth.NewAuthService(storeService)
	renderService := render.NewRenderService()

	characterService := service.NewCharacterService(storeService)
	identityPkg := &identity.Identity{}

	server := NewServer(authService, sessionService, identityPkg, characterService, renderService)

	enrichFunc := identity.EnrichContextWithUserID
	m := middleware.NewMiddlewareService(sessionService, enrichFunc)
	serverWithMiddleware := m.AuthMiddleware(server)

	slog.Info("server running", "on", "http://localhost:5050")
	log.Fatal(http.ListenAndServe(":5050", serverWithMiddleware))
}
