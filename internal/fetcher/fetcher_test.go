package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchPosts(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"userId": 1,
				"id": 1,
				"title": "Test Title",
				"body": "Test Body"
			}
		]`))
	}))
	defer server.Close()

	// Create a fetcher with the test server URL
	f := New(server.URL)

	// Test fetching posts
	posts, err := f.FetchPosts(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("Expected 1 post, got %d", len(posts))
	}

	if posts[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", posts[0].ID)
	}

	if posts[0].Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", posts[0].Title)
	}
}

func TestFetchPostsError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create a fetcher with the test server URL
	f := New(server.URL)

	// Test fetching posts with error
	_, err := f.FetchPosts(context.Background())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestFetchPostsInvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	// Create a fetcher with the test server URL
	f := New(server.URL)

	// Test fetching posts with invalid JSON
	_, err := f.FetchPosts(context.Background())
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}
}
