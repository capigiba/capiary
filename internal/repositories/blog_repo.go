package repositories

import (
	"context"
	"fmt"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/infra/db/mongodb"
	"github.com/capigiba/capiary/internal/infra/db/query"
	"go.mongodb.org/mongo-driver/bson"
)

type BlogPostRepository interface {
	Add(ctx context.Context, post entity.BlogPost) (string, error)
	UpdateByQuery(ctx context.Context, filter bson.M, update entity.BlogPost) error
	FindByQuery(ctx context.Context, opts query.QueryOptions) ([]entity.BlogPost, error)
	LoadAll(ctx context.Context) ([]entity.BlogPost, error)
}

type blogPostRepository struct {
	adapter *mongodb.MongoDBAdapter[entity.BlogPost]
}

func NewBlogPostRepository(db *mongodb.MongoDBClient) BlogPostRepository {
	return &blogPostRepository{
		adapter: mongodb.NewMongoDBAdapter[entity.BlogPost](
			db.GetClient(),
			"capiary",
			"blog",
		),
	}
}

func (r *blogPostRepository) Add(ctx context.Context, post entity.BlogPost) (string, error) {
	oid, err := r.adapter.InsertOne(post)
	if err != nil {
		return "", err
	}
	return oid.Hex(), nil
}

func (r *blogPostRepository) UpdateByQuery(ctx context.Context, filter bson.M, update entity.BlogPost) error {
	return r.adapter.UpdateOne(filter, update)
}

func (r *blogPostRepository) FindByQuery(ctx context.Context, opts query.QueryOptions) ([]entity.BlogPost, error) {
	fmt.Println(opts)
	return r.adapter.FindWithQuery(opts)
}

func (r *blogPostRepository) LoadAll(ctx context.Context) ([]entity.BlogPost, error) {
	return r.adapter.Find(bson.M{})
}
