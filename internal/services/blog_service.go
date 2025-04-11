package services

import (
	"context"
	"fmt"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/infra/db/query"
	"github.com/capigiba/capiary/internal/repositories"
)

type BlogPostService interface {
	CreatePost(ctx context.Context, post entity.BlogPost) (string, error)

	// This method will parse the raw filters/sorts in the service, build the QueryOptions, and then fetch.
	FindPostsWithRawQuery(ctx context.Context, rawFilters, rawSorts []string, rawFields string) ([]entity.BlogPost, error)

	// For updates/deletes, we can parse the filters in the service as well.
	UpdatePostByRawFilter(ctx context.Context, rawFilters []string, update entity.BlogPost) error
}

type blogPostService struct {
	repo repositories.BlogPostRepository
}

func NewBlogPostService(repo repositories.BlogPostRepository) BlogPostService {
	return &blogPostService{
		repo: repo,
	}
}

func (s *blogPostService) CreatePost(ctx context.Context, post entity.BlogPost) (string, error) {
	if post.Title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	return s.repo.Add(ctx, post)
}

// FindPostsWithRawQuery: parse the raw query params in the service, then build QueryOptions.
func (s *blogPostService) FindPostsWithRawQuery(ctx context.Context, rawFilters, rawSorts []string, rawFields string) ([]entity.BlogPost, error) {
	// 1) Parse filters
	parsedFilters, err := query.ParseFilters(rawFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	// 2) Parse sorts
	parsedSorts, err := query.ParseSorts(rawSorts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sorts: %w", err)
	}

	// 3) Parse fields
	parsedFields := query.ParseFields(rawFields)

	// 4) Build the QueryOptions
	opts := query.QueryOptions{
		Filters: parsedFilters,
		Sorts:   parsedSorts,
		Fields:  parsedFields,
	}

	// 5) Call the repository
	return s.repo.FindByQuery(ctx, opts)
}

// For update, we parse raw filters and build a bson.M filter. Then we call the repo method.
func (s *blogPostService) UpdatePostByRawFilter(ctx context.Context, rawFilters []string, update entity.BlogPost) error {
	parsedFilters, err := query.ParseFilters(rawFilters)
	if err != nil {
		return fmt.Errorf("failed to parse filters: %w", err)
	}

	filterDoc, _ := query.BuildMongoQuery(query.QueryOptions{
		Filters: parsedFilters})

	if update.Title == "" {
		return fmt.Errorf("cannot update post with empty title")
	}
	return s.repo.UpdateByQuery(ctx, filterDoc, update)
}
