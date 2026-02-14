package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"feedback/internal/config"
	"feedback/internal/db"
	"feedback/internal/modules/auth"
	"feedback/internal/modules/feedback"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize database pool
	ctx := context.Background()
	pool, err := db.InitPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer pool.Close()

	log.Println("Database connection established")

	// Create HTTP mux
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Register auth routes (Mailgun mail config)
	mailConfig := auth.MailConfig{
		APIKey:  cfg.MailgunAPIKey,
		Domain:  cfg.MailgunDomain,
		BaseURL: cfg.MailgunBaseURL,
		From:    cfg.EmailFrom,
	}
	auth.RegisterRoutes(mux, pool, cfg.JWTSecret, cfg.AppDeeplinkURL, mailConfig)

	// Register feedback routes
	feedback.RegisterRoutes(mux, pool, cfg.JWTSecret)

	// Create server (Render provides PORT as string)
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
