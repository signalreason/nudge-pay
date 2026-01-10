package config

import (
	"os"
	"strconv"
)

type Config struct {
	Addr          string
	DBPath        string
	JWTSecret     string
	WorkerEnabled bool
	BaseURL       string
}

func Load() Config {
	cfg := Config{
		Addr:          envOr("NUDGEPAY_ADDR", ":8080"),
		DBPath:        envOr("NUDGEPAY_DB", "./nudgepay.db"),
		JWTSecret:     envOr("NUDGEPAY_JWT_SECRET", "change-me"),
		WorkerEnabled: envBool("NUDGEPAY_WORKER", true),
		BaseURL:       envOr("NUDGEPAY_BASE_URL", "http://localhost:8080"),
	}
	return cfg
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}
