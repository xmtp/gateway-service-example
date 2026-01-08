package config

import (
	"log"
	"os"
	"strings"
)

type AuthMode string

const (
	AuthAllowAll AuthMode = "allowall"
	AuthJWTES256 AuthMode = "jwt_es256"
)

type Config struct {
	AuthMode AuthMode

	JWT struct {
		ExpectedIssuer string
		PublicKeyPEM   string
	}

	RateLimit struct {
		Enabled bool
	}
}

func Load() Config {
	mode := AuthMode(strings.ToLower(getEnv("GATEWAY_AUTH_MODE", "allowall")))

	cfg := Config{
		AuthMode: mode,
	}

	if mode == AuthJWTES256 {
		cfg.JWT.ExpectedIssuer = mustEnv("JWT_EXPECTED_ISSUER")
		cfg.JWT.PublicKeyPEM = mustEnv("JWT_PUBLIC_KEY")
	}

	cfg.RateLimit.Enabled = strings.ToLower(getEnv("RATE_LIMIT_ENABLED", "false")) == "true"

	return cfg
}

func mustEnv(key string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
