package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	BackendURL    string
	Provider      string
	LoadStorePath string
	Turvo         TurvoConfig
}

type TurvoConfig struct {
	BaseURL      string
	APIKey       string
	ClientName   string
	ClientSecret string
	Username     string
	Password     string
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
		Port:          getenvDefault("PORT", "8080"),
		BackendURL:    getenvDefault("BACKEND_URL", "http://localhost"),
		Provider:      getenvDefault("PROVIDER", "turvo"),
		LoadStorePath: getenvDefault("LOAD_STORE_PATH", ".data/drumkit-loads.json"),
		Turvo: TurvoConfig{
			BaseURL:      getenvDefault("TURVO_BASE_URL", "https://my-sandbox-publicapi.turvo.com/v1"),
			APIKey:       os.Getenv("TURVO_API"),
			ClientName:   os.Getenv("TURVO_CLIENT_NAME"),
			ClientSecret: os.Getenv("TURVO_CLIENT_SECRET"),
			Username:     os.Getenv("TURVO_USERNAME"),
			Password:     os.Getenv("TURVO_PASSWORD"),
		},
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

	switch strings.ToLower(strings.TrimSpace(c.Provider)) {
	case "", "turvo":
		for key, value := range map[string]string{
			"TURVO_API":           c.Turvo.APIKey,
			"TURVO_CLIENT_NAME":   c.Turvo.ClientName,
			"TURVO_CLIENT_SECRET": c.Turvo.ClientSecret,
			"TURVO_USERNAME":      c.Turvo.Username,
			"TURVO_PASSWORD":      c.Turvo.Password,
		} {
			if strings.TrimSpace(value) == "" {
				missing = append(missing, key)
			}
		}
	default:
		return fmt.Errorf("unsupported provider: %s", c.Provider)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}
