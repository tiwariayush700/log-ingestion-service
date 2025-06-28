package tracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestTracker(t *testing.T) (*Tracker, func()) {
	// Use a unique database name for each test to avoid conflicts
	dbName := "test_db_" + primitive.NewObjectID().Hex()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	tracker := &Tracker{
		client:     client,
		database:   dbName,
		collection: "ingest_status",
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

	return tracker, cleanup
}

func TestRecordSuccess(t *testing.T) {
	// Skip if no MongoDB connection
	if testing.Short() {
		t.Skip("Skipping MongoDB test in short mode")
	}

	tracker, cleanup := setupTestTracker(t)
	defer cleanup()

	// Record a successful ingestion
	ctx := context.Background()
	count := 10
	err := tracker.RecordSuccess(ctx, count)
	if err != nil {
		t.Fatalf("Failed to record success: %v", err)
	}

	// Retrieve the status
	status, err := tracker.GetLatestStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get latest status: %v", err)
	}

	// Verify results
	if !status.Success {
		t.Error("Expected success to be true")
	}

	if status.Count != count {
		t.Errorf("Expected count %d, got %d", count, status.Count)
	}

	if !status.Timestamp.Before(time.Now().UTC().Add(time.Second)) {
		t.Error("Expected timestamp to be in the past")
	}
}

func TestRecordFailure(t *testing.T) {
	// Skip if no MongoDB connection
	if testing.Short() {
		t.Skip("Skipping MongoDB test in short mode")
	}

	tracker, cleanup := setupTestTracker(t)
	defer cleanup()

	// Record a failed ingestion
	ctx := context.Background()
	testErr := errors.New("test error")
	err := tracker.RecordFailure(ctx, testErr)
	if err != nil {
		t.Fatalf("Failed to record failure: %v", err)
	}

	// Retrieve the status
	status, err := tracker.GetLatestStatus(ctx)
	if err != nil {
		t.Fatalf("Failed to get latest status: %v", err)
	}

	// Verify results
	if status.Success {
		t.Error("Expected success to be false")
	}

	if status.Error != testErr.Error() {
		t.Errorf("Expected error '%s', got '%s'", testErr.Error(), status.Error)
	}

	if !status.Timestamp.Before(time.Now().UTC().Add(time.Second)) {
		t.Error("Expected timestamp to be in the past")
	}
}

func TestGetLatestStatusNoRecords(t *testing.T) {
	// Skip if no MongoDB connection
	if testing.Short() {
		t.Skip("Skipping MongoDB test in short mode")
	}

	tracker, cleanup := setupTestTracker(t)
	defer cleanup()

	// Try to retrieve status when no records exist
	ctx := context.Background()
	_, err := tracker.GetLatestStatus(ctx)

	// Verify error
	if err == nil {
		t.Error("Expected error for no records, got nil")
	}
}
