package main

import (
	"database/sql"
	"fmt"

	"github.com/mp40/go-htmx-pccs/state"
)

func connectToStore(path string) (*sql.DB, error) {
	storeDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("init store db: %w", err)
	}

	_, err = storeDB.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY NOT NULL, email TEXT UNIQUE NOT NULL, hash TEXT NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL)")
	if err != nil {
		return nil, fmt.Errorf("db create user table: %w", err)
	}
	_, err = storeDB.Exec("CREATE TABLE IF NOT EXISTS characters (id TEXT PRIMARY KEY NOT NULL, user_id TEXT NOT NULL, name TEXT NOT NULL, str INTEGER NOT NULL, int INTEGER NOT NULL, wil INTEGER NOT NULL, hlt INTEGER NOT NULL, agi INTEGER NOT NULL, tch INTEGER NOT NULL, gun_combat_learning_points REAL NOT NULL, hand_to_hand_learning_points REAL NOT NULL, created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL)")
	if err != nil {
		return nil, fmt.Errorf("db create character table: %w", err)
	}

	return storeDB, nil
}

func connectToSession() (*sql.DB, error) {
	// found some gotchas with in memory - each connection gets private in mem db
	// maybe move out of memory (store or state db)
	sessionDB, err := sql.Open("sqlite", "file:sessions?mode=memory&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("init session db: %w", err)
	}
	// force one connection for now
	sessionDB.SetMaxOpenConns(1)
	_, err = sessionDB.Exec("CREATE TABLE IF NOT EXISTS sessions (id TEXT PRIMARY KEY NOT NULL, user_id TEXT NOT NULL, created_at INTEGER NOT NULL, expires_at INTEGER NOT NULL)")
	if err != nil {
		return nil, fmt.Errorf("db create session table: %w", err)
	}
	return sessionDB, nil
}

func connectToState(env string) (*sql.DB, error) {
	switch env {
	case "production":
		stateDB, err := sql.Open("sqlite", "/tmp/pccs_state.db")
		if err != nil {
			return nil, fmt.Errorf("init state db: %w", err)
		}
		if err := state.Bootstrap(stateDB); err != nil {
			return nil, fmt.Errorf("bootstrap state db: %w", err)
		}
		return stateDB, nil

	case "development", "":
		stateDB, err := sql.Open("sqlite", "./pccs_state.db")
		if err != nil {
			return nil, fmt.Errorf("init state db: %w", err)
		}
		err = state.Bootstrap(stateDB)
		if err != nil {
			return nil, fmt.Errorf("bootstrap state db: %w", err)
		}
		return stateDB, nil
	default:
		return nil, fmt.Errorf("connectToState: unknown env %q", env)
	}
}
