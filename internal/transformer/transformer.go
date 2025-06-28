package transformer

import (
	"time"

	"github.com/tiwariayush700/log-ingestion-service/internal/models"
)

// Transformer is responsible for transforming data
type Transformer struct {
	sourceName string
}

// New creates a new Transformer instance
func New(sourceName string) *Transformer {
	return &Transformer{
		sourceName: sourceName,
	}
}

// TransformPosts transforms posts by adding metadata
func (t *Transformer) TransformPosts(posts []models.Post) []models.EnrichedPost {
	enrichedPosts := make([]models.EnrichedPost, len(posts))
	now := time.Now().UTC()

	for i, post := range posts {
		enrichedPosts[i] = models.EnrichedPost{
			UserID:     post.UserID,
			PostID:     post.ID,
			Title:      post.Title,
			Body:       post.Body,
			IngestedAt: now,
			Source:     t.sourceName,
		}
	}

	return enrichedPosts
}
