package main

import (
	"log"
	"net/http"
	"time"

	"drumkit-take-home/internal/config"
	"drumkit-take-home/internal/drumkitstore"
	"drumkit-take-home/internal/httpapi"
	"drumkit-take-home/internal/turvo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	turvoClient := turvo.NewClient(cfg)
	loadStore, err := drumkitstore.New(cfg.LoadStorePath)
	if err != nil {
		log.Fatal(err)
	}
	server := httpapi.NewServer(turvoClient, loadStore)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("backend listening on %s", cfg.BackendURLWithPort())
	log.Fatal(httpServer.ListenAndServe())
}
