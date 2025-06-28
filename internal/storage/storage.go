package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/log-ingestion-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage handles data persistence
type Storage struct {
	client     *mongo.Client
	database   string
	collection string
}

// New creates a new Storage instance
func New(uri, database, collection string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Storage{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

// Close closes the database connection
func (s *Storage) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

// StorePosts stores the enriched posts in the database
func (s *Storage) StorePosts(ctx context.Context, posts []models.EnrichedPost) error {
	if len(posts) == 0 {
		return nil
	}

	collection := s.client.Database(s.database).Collection(s.collection)

	// Convert posts to interface slice for bulk write
	documents := make([]interface{}, len(posts))
	for i, post := range posts {
		documents[i] = post
	}

	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to insert posts: %w", err)
	}

	return nil
}

// GetPosts retrieves all posts from the database
func (s *Storage) GetPosts(ctx context.Context) ([]models.EnrichedPost, error) {
	collection := s.client.Database(s.database).Collection(s.collection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []models.EnrichedPost
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}

	return posts, nil
}

// GetPostByID retrieves a post by its ID
func (s *Storage) GetPostByID(ctx context.Context, id string) (models.EnrichedPost, error) {
	collection := s.client.Database(s.database).Collection(s.collection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.EnrichedPost{}, fmt.Errorf("invalid ID format: %w", err)
	}

	var post models.EnrichedPost
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		return models.EnrichedPost{}, fmt.Errorf("failed to find post: %w", err)
	}

	return post, nil
}
