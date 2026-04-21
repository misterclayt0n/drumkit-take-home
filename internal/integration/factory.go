package integration

import (
	"fmt"
	"strings"

	"drumkit-take-home/internal/config"
	"drumkit-take-home/internal/turvo"
)

func NewProvider(cfg config.Config) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case "", "turvo":
		return turvo.NewClient(cfg.Turvo), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
