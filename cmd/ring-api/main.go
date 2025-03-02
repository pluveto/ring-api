package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pluveto/ring-api/api"
	"github.com/pluveto/ring-api/internal/player"
	"github.com/pluveto/ring-api/pkg/config"
)

func main() {
	cfg := config.Load()
	p := player.New()
	handler := api.NewHandler(p, cfg.AudioPath)

	http.HandleFunc("/api/ring", handler.HandleRing)

	log.Printf("Starting server on :%s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed: %v\n", err)
		os.Exit(1)
	}
}
