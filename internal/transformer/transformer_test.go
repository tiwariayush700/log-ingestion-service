package transformer

import (
	"testing"
	"time"

	"github.com/tiwariayush700/log-ingestion-service/internal/models"
)

func TestTransformPosts(t *testing.T) {
	// Create a transformer
	sourceName := "test_source"
	transformer := New(sourceName)

	// Create test posts
	posts := []models.Post{
		{
			UserID: 1,
			ID:     1,
			Title:  "Test Title",
			Body:   "Test Body",
		},
		{
			UserID: 2,
			ID:     2,
			Title:  "Another Title",
			Body:   "Another Body",
		},
	}

	// Transform posts
	before := time.Now().UTC()
	enrichedPosts := transformer.TransformPosts(posts)
	after := time.Now().UTC()

	// Check results
	if len(enrichedPosts) != len(posts) {
		t.Fatalf("Expected %d enriched posts, got %d", len(posts), len(enrichedPosts))
	}

	for i, post := range posts {
		enriched := enrichedPosts[i]

		if enriched.UserID != post.UserID {
			t.Errorf("Expected UserID %d, got %d", post.UserID, enriched.UserID)
		}

		if enriched.PostID != post.ID {
			t.Errorf("Expected PostID %d, got %d", post.ID, enriched.PostID)
		}

		if enriched.Title != post.Title {
			t.Errorf("Expected Title '%s', got '%s'", post.Title, enriched.Title)
		}

		if enriched.Body != post.Body {
			t.Errorf("Expected Body '%s', got '%s'", post.Body, enriched.Body)
		}

		if enriched.Source != sourceName {
			t.Errorf("Expected Source '%s', got '%s'", sourceName, enriched.Source)
		}

		// Check that IngestedAt is between before and after
		if enriched.IngestedAt.Before(before) || enriched.IngestedAt.After(after) {
			t.Errorf("IngestedAt timestamp %v is not between %v and %v", enriched.IngestedAt, before, after)
		}
	}
}

func TestTransformEmptyPosts(t *testing.T) {
	// Create a transformer
	transformer := New("test_source")

	// Transform empty posts
	enrichedPosts := transformer.TransformPosts([]models.Post{})

	// Check results
	if len(enrichedPosts) != 0 {
		t.Fatalf("Expected 0 enriched posts, got %d", len(enrichedPosts))
	}
}
