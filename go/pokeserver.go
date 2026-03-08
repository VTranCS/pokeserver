// Package main implements a Pokemon voting server with clean architecture.
//
// The server provides REST API endpoints for:
//   - Getting random Pokemon from PokeAPI
//   - Voting on Pokemon (up/down votes)  
//   - Retrieving all Pokemon with vote counts
//   - Health checks
//
// The application follows clean architecture principles with proper
// separation of concerns, dependency injection, and structured logging.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var repo *Repository

func main() {
	// Initialize logger
	logger := NewLogger()
	logger.Info("Starting pokeserver application")

	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to load configuration")
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Initialize database repository
	ctx := context.Background()
	repo, err = NewRepository(ctx, config.Database.URL, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to database")
		log.Fatalf("Error connecting to database: %s", err)
	}
	defer repo.Close()

	// Initialize Pokemon client
	pokeClient := NewPokeClient(config.PokeAPI.URL, config.PokeAPI.Max, logger)

	// Initialize Pokemon service
	pokemonService := NewPokemonService(repo, pokeClient, logger)

	// Initialize database
	if err := pokemonService.InitializeDatabase(ctx); err != nil {
		logger.WithError(err).Error("Failed to initialize database")
		log.Fatalf("Error initializing database: %s", err)
	}

	// Create HTTP handlers
	handlers := NewHandlers(pokemonService, logger)

	// Setup middleware
	middleware := ChainMiddleware(
		RecoveryMiddleware(logger),
		LoggingMiddleware(logger),
		CORSMiddleware(),
	)

	// Setup routes
	mux := http.NewServeMux()
	mux.Handle("/getall", middleware(http.HandlerFunc(handlers.HandleGetAllPokemon)))
	mux.Handle("/getpokemon", middleware(http.HandlerFunc(handlers.HandleGetPokemon)))
	mux.Handle("/vote", middleware(http.HandlerFunc(handlers.HandleVote)))
	mux.Handle("/health", middleware(http.HandlerFunc(handlers.HandleHealth)))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/", middleware(http.HandlerFunc(handlers.HandlePokeStop)))

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		logger.WithFields(map[string]any{
			"port": config.Server.Port,
		}).Info("Starting HTTP server")
		
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("HTTP server failed to start")
			log.Fatalf("HTTP server failed to start: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give the server a timeout to finish handling requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	logger.Info("Server exited")
}

// PokeDBEntry represents a Pokemon entry in the database with vote information.
type PokeDBEntry struct {
	Id   int    `json:"id"`   // Pokemon ID from PokeAPI
	Name string `json:"name"` // Pokemon name
	Vote int    `json:"vote"` // Current vote count
	Url  string `json:"url"`  // URL to Pokemon sprite image
}
