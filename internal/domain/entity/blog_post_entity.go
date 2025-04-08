package entity

import "time"

// BlogPost describes the top-level information for a blog post.
type BlogPost struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Blocks    []Block   `json:"blocks"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
