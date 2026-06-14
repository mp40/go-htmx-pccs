package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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
	"github.com/mp40/go-htmx-pccs/state"
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
	config, err := loadConfig(getenv)
	if err != nil {
		return err
	}

	storeDB, err := connectToStore(config.StoreDBPath)
	if err != nil {
		return err
	}
	defer storeDB.Close()

	sessionDB, err := connectToSession()
	if err != nil {
		return err
	}
	defer sessionDB.Close()

	stateDB, err := connectToState(config.Env)
	if err != nil {
		return err
	}
	defer stateDB.Close()

	storeService := store.NewStoreService(storeDB)
	stateService := state.NewStateService(stateDB)
	sessionService := session.NewSessionService(sessionDB)
	authService := auth.NewAuthService(storeService, config.BcryptCost)
	renderService := render.NewRenderService()

	characterService := service.NewCharacterService(storeService)
	gearService := service.NewGearService(stateService)
	identityPkg := &identity.Identity{}

	s := server.NewServer(authService, sessionService, identityPkg, characterService, gearService, renderService)

	enrichFunc := identity.EnrichContextWithUserID
	m := middleware.NewMiddlewareService(sessionService, enrichFunc)
	serverWithMiddleware := m.AuthMiddleware(s)

	httpServer := &http.Server{
		Addr:         ":" + config.Port,
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
