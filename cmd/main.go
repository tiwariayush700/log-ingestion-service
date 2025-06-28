package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yourusername/log-ingestion-service/config"
	"github.com/yourusername/log-ingestion-service/internal/api"
	"github.com/yourusername/log-ingestion-service/internal/fetcher"
	"github.com/yourusername/log-ingestion-service/internal/storage"
	"github.com/yourusername/log-ingestion-service/internal/tracker"
	"github.com/yourusername/log-ingestion-service/internal/transformer"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize components
	fetch := fetcher.New(cfg.APIEndpoint)
	transform := transformer.New(cfg.SourceName)

	store, err := storage.New(cfg.MongoURI, cfg.MongoDatabase, cfg.MongoCollection)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	track, err := tracker.New(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("Failed to initialize tracker: %v", err)
	}

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the API server
	apiServer := api.New(store, track)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := apiServer.Run(":" + cfg.ServerPort); err != nil {
			log.Printf("API server error: %v", err)
			cancel()
		}
	}()

	// Start the ingestion process
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(cfg.FetchInterval)
		defer ticker.Stop()

		// Run immediately on startup
		ingestData(ctx, fetch, transform, store, track)

		for {
			select {
			case <-ticker.C:
				ingestData(ctx, fetch, transform, store, track)
			case <-ctx.Done():
				return
			}
		}
	}()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
		log.Println("Received shutdown signal")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	// Clean up resources
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := store.Close(shutdownCtx); err != nil {
		log.Printf("Error closing storage: %v", err)
	}

	if err := track.Close(shutdownCtx); err != nil {
		log.Printf("Error closing tracker: %v", err)
	}

	cancel()  // Cancel the context to stop all goroutines
	wg.Wait() // Wait for all goroutines to finish
	log.Println("Application shutdown complete")
}

func ingestData(ctx context.Context, fetch *fetcher.Fetcher, transform *transformer.Transformer, store *storage.Storage, track *tracker.Tracker) {
	log.Println("Starting data ingestion...")

	// Fetch data
	posts, err := fetch.FetchPosts(ctx)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		if trackErr := track.RecordFailure(ctx, err); trackErr != nil {
			log.Printf("Error recording failure: %v", trackErr)
		}
		return
	}

	// Transform data
	enrichedPosts := transform.TransformPosts(posts)

	// Store data
	if err := store.StorePosts(ctx, enrichedPosts); err != nil {
		log.Printf("Error storing posts: %v", err)
		if trackErr := track.RecordFailure(ctx, err); trackErr != nil {
			log.Printf("Error recording failure: %v", trackErr)
		}
		return
	}

	// Record success
	if err := track.RecordSuccess(ctx, len(enrichedPosts)); err != nil {
		log.Printf("Error recording success: %v", err)
	}

	log.Printf("Successfully ingested %d posts", len(enrichedPosts))
}
