package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogPost describes the top-level information for a blog post.
type BlogPost struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title"`
	Blocks    []Block            `json:"blocks"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
