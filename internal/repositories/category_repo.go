package repositories

import (
	"context"
	"fmt"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/infra/db/mongodb"
	"github.com/capigiba/capiary/internal/infra/db/query"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryRepository interface {
	Add(ctx context.Context, post entity.Category) (string, error)
	UpdateByQuery(ctx context.Context, filter bson.M, update entity.Category) error
	FindByQuery(ctx context.Context, opts query.QueryOptions) ([]entity.Category, error)
	LoadAll(ctx context.Context) ([]entity.Category, error)
}

type categoryRepository struct {
	adapter *mongodb.MongoDBAdapter[entity.Category]
}

func NewCategoryRepository(db *mongodb.MongoDBClient) CategoryRepository {
	return &categoryRepository{
		adapter: mongodb.NewMongoDBAdapter[entity.Category](
			db.GetClient(),
			"capiary",
			"category",
		),
	}
}

func (r *categoryRepository) Add(ctx context.Context, post entity.Category) (string, error) {
	oid, err := r.adapter.InsertOne(post)
	if err != nil {
		return "", err
	}
	return oid.Hex(), nil
}

func (r *categoryRepository) UpdateByQuery(ctx context.Context, filter bson.M, update entity.Category) error {
	return r.adapter.UpdateOne(filter, update)
}

func (r *categoryRepository) FindByQuery(ctx context.Context, opts query.QueryOptions) ([]entity.Category, error) {
	fmt.Println(opts)
	return r.adapter.FindWithQuery(opts)
}

func (r *categoryRepository) LoadAll(ctx context.Context) ([]entity.Category, error) {
	loadAllOpts := query.QueryOptions{
		Sorts: []query.Sort{
			{
				Field: "name",
				Desc:  false,
			},
		},
	}
	return r.adapter.FindWithQuery(loadAllOpts)
}
