package main

import (
	"log"
	"net/http"
	"time"

	"drumkit-take-home/internal/config"
	"drumkit-take-home/internal/drumkitstore"
	"drumkit-take-home/internal/httpapi"
	"drumkit-take-home/internal/integration"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	provider, err := integration.NewProvider(cfg)
	if err != nil {
		log.Fatal(err)
	}

	loadStore, err := drumkitstore.New(cfg.LoadStorePath)
	if err != nil {
		log.Fatal(err)
	}
	server := httpapi.NewServer(provider, loadStore)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("backend listening on %s", cfg.BackendURLWithPort())
	log.Fatal(httpServer.ListenAndServe())
}
