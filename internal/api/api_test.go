package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/log-ingestion-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockStorage is a mock implementation of the storage interface
type MockStorage struct {
	posts []models.EnrichedPost
}

func (m *MockStorage) GetPosts(ctx interface{}) ([]models.EnrichedPost, error) {
	return m.posts, nil
}

func (m *MockStorage) GetPostByID(ctx interface{}, id string) (models.EnrichedPost, error) {
	for _, post := range m.posts {
		if post.ID.Hex() == id {
			return post, nil
		}
	}
	return models.EnrichedPost{}, nil
}

// MockTracker is a mock implementation of the tracker interface
type MockTracker struct {
	status models.IngestStatus
}

func (m *MockTracker) GetLatestStatus(ctx interface{}) (models.IngestStatus, error) {
	return m.status, nil
}

func setupTestAPI() (*API, *MockStorage, *MockTracker) {
	gin.SetMode(gin.TestMode)

	// Create mock storage with test data
	objID, _ := primitive.ObjectIDFromHex("5f50c31f5dc4b6d5c8456e77")
	mockStorage := &MockStorage{
		posts: []models.EnrichedPost{
			{
				ID:         objID,
				UserID:     1,
				PostID:     1,
				Title:      "Test Title",
				Body:       "Test Body",
				IngestedAt: time.Now().UTC(),
				Source:     "test_source",
			},
		},
	}

	// Create mock tracker with test data
	mockTracker := &MockTracker{
		status: models.IngestStatus{
			ID:        primitive.NewObjectID(),
			Timestamp: time.Now().UTC(),
			Success:   true,
			Count:     1,
		},
	}

	// Create API with mock dependencies
	api := &API{
		router:  gin.New(),
		storage: mockStorage,
		tracker: mockTracker,
	}
	api.setupRoutes()

	return api, mockStorage, mockTracker
}

func TestGetLogs(t *testing.T) {
	api, _, _ := setupTestAPI()

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/api/logs", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	api.router.ServeHTTP(resp, req)

	// Check response
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var logs []models.EnrichedPost
	if err := json.Unmarshal(resp.Body.Bytes(), &logs); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if logs[0].Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", logs[0].Title)
	}
}

func TestGetLogByID(t *testing.T) {
	api, _, _ := setupTestAPI()

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/api/logs/5f50c31f5dc4b6d5c8456e77", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	api.router.ServeHTTP(resp, req)

	// Check response
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var log models.EnrichedPost
	if err := json.Unmarshal(resp.Body.Bytes(), &log); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if log.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", log.Title)
	}
}

func TestGetStatus(t *testing.T) {
	api, _, _ := setupTestAPI()

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	api.router.ServeHTTP(resp, req)

	// Check response
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var status models.IngestStatus
	if err := json.Unmarshal(resp.Body.Bytes(), &status); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !status.Success {
		t.Error("Expected success to be true")
	}

	if status.Count != 1 {
		t.Errorf("Expected count 1, got %d", status.Count)
	}
}
