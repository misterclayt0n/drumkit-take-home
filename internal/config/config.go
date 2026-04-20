package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	BackendURL        string
	TurvoBaseURL      string
	TurvoAPIKey       string
	TurvoClientName   string
	TurvoClientSecret string
	TurvoUsername     string
	TurvoPassword     string
	LoadStorePath     string
}

func Load() (Config, error) {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	getenvDefault := func(key, fallback string) string {
		if value := os.Getenv(key); strings.TrimSpace(value) != "" {
			return value
		}

		return fallback
	}

	cfg := Config{
		Port:              getenvDefault("PORT", "8080"),
		BackendURL:        getenvDefault("BACKEND_URL", "http://localhost"),
		TurvoBaseURL:      getenvDefault("TURVO_BASE_URL", "https://my-sandbox-publicapi.turvo.com/v1"),
		TurvoAPIKey:       os.Getenv("TURVO_API"),
		TurvoClientName:   os.Getenv("TURVO_CLIENT_NAME"),
		TurvoClientSecret: os.Getenv("TURVO_CLIENT_SECRET"),
		TurvoUsername:     os.Getenv("TURVO_USERNAME"),
		TurvoPassword:     os.Getenv("TURVO_PASSWORD"),
		LoadStorePath:     getenvDefault("LOAD_STORE_PATH", ".data/drumkit-loads.json"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) BackendURLWithPort() string {
	return fmt.Sprintf("%s:%s", strings.TrimSuffix(c.BackendURL, "/"), c.Port)
}

func (c Config) Validate() error {
	var missing []string

	for key, value := range map[string]string{
		"TURVO_API":           c.TurvoAPIKey,
		"TURVO_CLIENT_NAME":   c.TurvoClientName,
		"TURVO_CLIENT_SECRET": c.TurvoClientSecret,
		"TURVO_USERNAME":      c.TurvoUsername,
		"TURVO_PASSWORD":      c.TurvoPassword,
	} {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}
