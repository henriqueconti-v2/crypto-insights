package main

import (
	"flag"
	"log"
	"net/http"

	"crypto-alerts/internal/config"
	"crypto-alerts/internal/handler"
)

const (
	defaultPort = "8080"
)

func main() {
	port := flag.String("port", defaultPort, "Port to run the server on")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	api, err := handler.NewAPI(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	mux := api.SetupRoutes()
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	log.Printf("Server running on port %s", *port)
	log.Printf("Database connection established to: %s/%s", cfg.Database.Host, cfg.Database.DBName)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
