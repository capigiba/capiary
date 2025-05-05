package entity

import (
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogPost describes the top-level information for a blog post.
type BlogPost struct {
	ID         primitive.ObjectID  `json:"id" bson:"_id"`
	AuthorID   int                 `json:"author_id"`
	Categories []string            `json:"categories"`
	Title      string              `json:"title"`
	Blocks     []Block             `json:"blocks"`
	Status     constant.BlogStatus `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}
