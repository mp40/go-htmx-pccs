package main

import (
	"fmt"
	"strconv"
)

type Config struct {
	Env         string
	Port        string
	BcryptCost  int
	StoreDBPath string
}

func loadConfig(getenv func(string) string) (Config, error) {
	config := Config{}
	costStr := getenv("COST")
	if costStr == "" {
		return config, fmt.Errorf("COST env var required")
	}
	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return config, fmt.Errorf("COST must be an integer: %w", err)
	}

	envStr := getenv("ENV")

	port := getenv("PORT")
	if port == "" {
		port = "5050"
	}

	path := getenv("STORE_PATH")
	if path == "" {
		path = "./pccs_store.db"
	}

	cfg := Config{
		Env:         envStr,
		Port:        port,
		StoreDBPath: path,
		BcryptCost:  cost,
	}

	return cfg, nil
}
