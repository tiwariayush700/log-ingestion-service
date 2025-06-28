package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/log-ingestion-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Tracker monitors ingestion progress
type Tracker struct {
	client     *mongo.Client
	database   string
	collection string
}

// New creates a new Tracker instance
func New(uri, database string) (*Tracker, error) {
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

	return &Tracker{
		client:     client,
		database:   database,
		collection: "ingest_status",
	}, nil
}

// Close closes the database connection
func (t *Tracker) Close(ctx context.Context) error {
	return t.client.Disconnect(ctx)
}

// RecordSuccess records a successful ingestion
func (t *Tracker) RecordSuccess(ctx context.Context, count int) error {
	collection := t.client.Database(t.database).Collection(t.collection)

	status := models.IngestStatus{
		Timestamp: time.Now().UTC(),
		Success:   true,
		Count:     count,
	}

	_, err := collection.InsertOne(ctx, status)
	if err != nil {
		return fmt.Errorf("failed to record success: %w", err)
	}

	return nil
}

// RecordFailure records a failed ingestion
func (t *Tracker) RecordFailure(ctx context.Context, err error) error {
	collection := t.client.Database(t.database).Collection(t.collection)

	status := models.IngestStatus{
		Timestamp: time.Now().UTC(),
		Success:   false,
		Error:     err.Error(),
	}

	_, err = collection.InsertOne(ctx, status)
	if err != nil {
		return fmt.Errorf("failed to record failure: %w", err)
	}

	return nil
}

// GetLatestStatus retrieves the latest ingestion status
func (t *Tracker) GetLatestStatus(ctx context.Context) (models.IngestStatus, error) {
	collection := t.client.Database(t.database).Collection(t.collection)

	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	var status models.IngestStatus
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&status)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.IngestStatus{}, fmt.Errorf("no ingestion status found")
		}
		return models.IngestStatus{}, fmt.Errorf("failed to get latest status: %w", err)
	}

	return status, nil
}
