package services

import (
	"context"
	"fmt"
	"time"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/infra/db/query"
	"github.com/capigiba/capiary/internal/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryService interface {
	Create(c *gin.Context, category entity.Category) (string, error)
	Find(ctx context.Context, rawFilters, rawSorts []string, rawFields string, page, pageSize int) ([]entity.Category, error)
	UpdateByRawFilter(c *gin.Context, rawFilters []string, update entity.Category) error
	LoadAll(ctx context.Context) ([]entity.Category, error)
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) Create(c *gin.Context, category entity.Category) (string, error) {
	if category.Name == "" {
		return "", fmt.Errorf("cannot update post with empty name")
	}
	insertedID, err := s.repo.Add(c.Request.Context(), category)
	if err != nil {
		return "", fmt.Errorf("failed to insert new category: %w", err)
	}

	return insertedID, nil
}

func (s *categoryService) Find(
	ctx context.Context,
	rawFilters, rawSorts []string,
	rawFields string,
	page, pageSize int,
) ([]entity.Category, error) {
	parsedFilters, err := query.ParseFilters(rawFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	// Convert "id" filter to "_id" if present:
	for i, f := range parsedFilters {
		if f.Field == "id" {
			idStr, ok := f.Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid id filter value type")
			}
			oid, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert id to ObjectID: %w", err)
			}
			parsedFilters[i].Field = "_id"
			parsedFilters[i].Value = oid
		}
	}

	parsedSorts, err := query.ParseSorts(rawSorts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sorts: %w", err)
	}
	if len(parsedSorts) == 0 {
		parsedSorts = []query.Sort{{Field: "created_at", Desc: true}}
	}
	parsedSorts = append(parsedSorts, query.Sort{Field: "_id", Desc: true})

	parsedFields := query.ParseFields(rawFields)

	opts := query.QueryOptions{
		Filters: parsedFilters,
		Sorts:   parsedSorts,
		Fields:  parsedFields,
		Skip:    int64((page - 1) * pageSize),
		Limit:   int64(pageSize),
	}

	categories, err := s.repo.FindByQuery(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) UpdateByRawFilter(c *gin.Context, rawFilters []string, update entity.Category) error {
	parsedFilters, err := query.ParseFilters(rawFilters)
	if err != nil {
		return fmt.Errorf("failed to parse filters: %w", err)

	}

	for i, f := range parsedFilters {
		if f.Field == "id" {
			if hex, ok := f.Value.(string); ok {
				if oid, err := primitive.ObjectIDFromHex(hex); err == nil {
					parsedFilters[i].Field = "_id"
					parsedFilters[i].Value = oid
				}
			}
		}
	}

	filterDoc, _ := query.BuildMongoQuery(query.QueryOptions{
		Filters: parsedFilters})

	if update.Name == "" {
		return fmt.Errorf("cannot update post with empty name")
	}
	update.UpdatedAt = time.Now()
	setDoc := bson.M{
		"name":        update.Name,
		"description": update.Description,
		"updated_at":  update.UpdatedAt,
	}
	return s.repo.UpdateFieldsByQuery(c.Request.Context(), filterDoc, setDoc)
}

func (s *categoryService) LoadAll(ctx context.Context) ([]entity.Category, error) {
	posts, err := s.repo.LoadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load all posts: %w", err)
	}
	return posts, nil
}
