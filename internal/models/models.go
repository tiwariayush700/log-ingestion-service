package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents the data structure from the JSONPlaceholder API
type Post struct {
	UserID int    `json:"userId" bson:"userId"`
	ID     int    `json:"id" bson:"id"`
	Title  string `json:"title" bson:"title"`
	Body   string `json:"body" bson:"body"`
}

// EnrichedPost represents a post with additional metadata
type EnrichedPost struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     int                `json:"userId" bson:"userId"`
	PostID     int                `json:"postId" bson:"postId"`
	Title      string             `json:"title" bson:"title"`
	Body       string             `json:"body" bson:"body"`
	IngestedAt time.Time          `json:"ingested_at" bson:"ingested_at"`
	Source     string             `json:"source" bson:"source"`
}

// IngestStatus represents the status of the latest ingestion
type IngestStatus struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Success   bool               `json:"success" bson:"success"`
	Count     int                `json:"count" bson:"count"`
	Error     string             `json:"error,omitempty" bson:"error,omitempty"`
}
