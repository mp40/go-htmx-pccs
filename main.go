package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mp40/go-htmx-pccs/auth"
	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/middleware"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/server"
	"github.com/mp40/go-htmx-pccs/service"
	"github.com/mp40/go-htmx-pccs/session"
	"github.com/mp40/go-htmx-pccs/store"
	_ "modernc.org/sqlite"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	_ = godotenv.Load()
	if err := run(ctx, os.Getenv, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// in future
// args []string could be good when want to build dev versions for testing ect
func run(ctx context.Context, getenv func(string) string, stderr io.Writer) error {
	costStr := getenv("COST")
	if costStr == "" {
		return fmt.Errorf("COST env var required")
	}
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return fmt.Errorf("COST must be an integer: %w", err)
	}

	storeDB, err := sql.Open("sqlite", "./pccs_store.db")
	if err != nil {
		return fmt.Errorf("init store db: %w", err)
	}
	defer storeDB.Close()
	// fine for now but need to do this only if in local dev mode once deployed
	_, err = storeDB.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY NOT NULL, email TEXT UNIQUE NOT NULL, hash TEXT NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL)")
	if err != nil {
		return fmt.Errorf("db create user table: %w", err)
	}
	// fine for now but need to do this only if in local dev mode once deployed
	_, err = storeDB.Exec("CREATE TABLE IF NOT EXISTS characters (id TEXT PRIMARY KEY NOT NULL, user_id TEXT NOT NULL, name TEXT NOT NULL, str INTEGER NOT NULL, int INTEGER NOT NULL, wil INTEGER NOT NULL, hlt INTEGER NOT NULL, agi INTEGER NOT NULL, tch INTEGER NOT NULL, gun_combat_learning_points REAL NOT NULL, hand_to_hand_learning_points REAL NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL)")
	if err != nil {
		return fmt.Errorf("db create character table: %w", err)
	}

	// found some gotchas with in memory - each connection gets private in mem db
	// maybe move out of memory (store or state db)
	sessionDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return fmt.Errorf("init session db: %w", err)
	}
	// force one connection for now
	sessionDB.SetMaxOpenConns(1)
	defer sessionDB.Close()
	// fine for now but need to do this only if in local dev mode once deployed
	_, err = sessionDB.Exec("CREATE TABLE IF NOT EXISTS sessions (id TEXT PRIMARY KEY NOT NULL, user_id TEXT NOT NULL, created_at INTEGER NOT NULL, expires_at INTEGER NOT NULL)")
	if err != nil {
		return fmt.Errorf("db create session table: %w", err)
	}

	// third SQLite db to be called state
	// a read only db ships with app for pccs game data

	storeService := store.NewStoreService(storeDB)
	sessionService := session.NewSessionService(sessionDB)
	authService := auth.NewAuthService(storeService, cost)
	renderService := render.NewRenderService()

	characterService := service.NewCharacterService(storeService)
	identityPkg := &identity.Identity{}

	s := server.NewServer(authService, sessionService, identityPkg, characterService, renderService)

	enrichFunc := identity.EnrichContextWithUserID
	m := middleware.NewMiddlewareService(sessionService, enrichFunc)
	serverWithMiddleware := m.AuthMiddleware(s)

	port := getenv("PORT")
	if port == "" {
		port = "5050"
	}

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      serverWithMiddleware,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		slog.Info("server running", "addr", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
			cancel()
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}
