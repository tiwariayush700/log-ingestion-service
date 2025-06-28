package storage

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/log-ingestion-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestStorage(t *testing.T) (*Storage, func()) {
	// Use a unique database name for each test to avoid conflicts
	dbName := "test_db_" + primitive.NewObjectID().Hex()
	collName := "test_collection"

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	storage := &Storage{
		client:     client,
		database:   dbName,
		collection: collName,
	}

	// Return a cleanup function
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Database(dbName).Drop(ctx); err != nil {
			t.Logf("Failed to drop test database: %v", err)
		}

		if err := client.Disconnect(ctx); err != nil {
			t.Logf("Failed to disconnect from MongoDB: %v", err)
		}
	}

	return storage, cleanup
}

func TestStorePosts(t *testing.T) {
	// Skip if no MongoDB connection
	if testing.Short() {
		t.Skip("Skipping MongoDB test in short mode")
	}

	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create test data
	posts := []models.EnrichedPost{
		{
			UserID:     1,
			PostID:     1,
			Title:      "Test Title 1",
			Body:       "Test Body 1",
			IngestedAt: time.Now().UTC(),
			Source:     "test_source",
		},
		{
			UserID:     2,
			PostID:     2,
			Title:      "Test Title 2",
			Body:       "Test Body 2",
			IngestedAt: time.Now().UTC(),
			Source:     "test_source",
		},
	}

	// Store posts
	ctx := context.Background()
	err := storage.StorePosts(ctx, posts)
	if err != nil {
		t.Fatalf("Failed to store posts: %v", err)
	}

	// Retrieve posts
	retrievedPosts, err := storage.GetPosts(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve posts: %v", err)
	}

	// Verify results
	if len(retrievedPosts) != len(posts) {
		t.Errorf("Expected %d posts, got %d", len(posts), len(retrievedPosts))
	}
}

func TestGetPostByID(t *testing.T) {
	// Skip if no MongoDB connection
	if testing.Short() {
		t.Skip("Skipping MongoDB test in short mode")
	}

	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Create and store a post
	post := models.EnrichedPost{
		UserID:     1,
		PostID:     1,
		Title:      "Test Title",
		Body:       "Test Body",
		IngestedAt: time.Now().UTC(),
		Source:     "test_source",
	}

	ctx := context.Background()
	collection := storage.client.Database(storage.database).Collection(storage.collection)
	result, err := collection.InsertOne(ctx, post)
	if err != nil {
		t.Fatalf("Failed to insert test post: %v", err)
	}

	// Get the inserted ID
	id := result.InsertedID.(primitive.ObjectID).Hex()

	// Retrieve the post by ID
	retrievedPost, err := storage.GetPostByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to retrieve post by ID: %v", err)
	}

	// Verify results
	if retrievedPost.Title != post.Title {
		t.Errorf("Expected title '%s', got '%s'", post.Title, retrievedPost.Title)
	}

	if retrievedPost.Body != post.Body {
		t.Errorf("Expected body '%s', got '%s'", post.Body, retrievedPost.Body)
	}
}

func TestGetPostByIDNotFound(t *testing.T) {
	// Skip if no MongoDB connection
	if testing.Short() {
		t.Skip("Skipping MongoDB test in short mode")
	}

	storage, cleanup := setupTestStorage(t)
	defer cleanup()

	// Try to retrieve a non-existent post
	ctx := context.Background()
	id := primitive.NewObjectID().Hex()
	_, err := storage.GetPostByID(ctx, id)

	// Verify error
	if err == nil {
		t.Error("Expected error for non-existent post, got nil")
	}
}
