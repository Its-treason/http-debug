package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"its-treason/web-test/db"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Configuration error", "error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.New(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to initialize database connections", "error", err)
	}
	defer database.Close()

	server := NewServer(database)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	go func() {
		log.Info("Starting HTTP server", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "error", err)
	}
}
