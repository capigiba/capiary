package services

import (
	"context"
	"fmt"
	"time"

	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/infra/db/query"
	"github.com/capigiba/capiary/internal/infra/storage"
	"github.com/capigiba/capiary/internal/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogPostService interface {
	CreatePostWithFiles(c *gin.Context, post entity.BlogPost) (string, error)
	FindPostsWithRawQuery(ctx context.Context, rawFilters, rawSorts []string, rawFields string, page, pageSize int) ([]entity.BlogPost, error)
	UpdatePostByRawFilter(ctx context.Context, rawFilters []string, update entity.BlogPost) error
	LoadAllPosts(ctx context.Context) ([]entity.BlogPost, error)
	SoftDeletePostByRawFilter(ctx context.Context, rawFilters []string) error
}

type blogPostService struct {
	repo       repositories.BlogPostRepository
	s3Uploader storage.S3UploaderInterface
}

func NewBlogPostService(repo repositories.BlogPostRepository, s3Uploader storage.S3UploaderInterface) BlogPostService {
	return &blogPostService{
		repo:       repo,
		s3Uploader: s3Uploader,
	}
}

func (s *blogPostService) CreatePostWithFiles(c *gin.Context, post entity.BlogPost) (string, error) {
	if post.Title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}

	// Loop over the blocks
	for i := range post.Blocks {
		switch post.Blocks[i].Type {
		case entity.BlockTypeImage:
			if post.Blocks[i].Image != nil {
				key := fmt.Sprintf("block_%d_fileBytes", i)
				raw, exists := c.Get(key)
				if !exists {
					return "", fmt.Errorf("missing file data for image block %d", i)
				}
				fileBytes, ok := raw.([]byte)
				if !ok {
					return "", fmt.Errorf("invalid file data format for image block %d", i)
				}

				// Call S3 upload
				s3Key, err := s.s3Uploader.UploadFile(
					constant.S3FolderImage,
					post.Blocks[i].Image.Filename,
					"image/png",
					"0",
					fileBytes,
				)
				if err != nil {
					return "", fmt.Errorf("failed to upload image block %d: %w", i, err)
				}

				// Overwrite the Filename with the returned S3 key
				post.Blocks[i].Image.Filename = s3Key
			}

		case entity.BlockTypeVideo:
			if post.Blocks[i].Video != nil {
				key := fmt.Sprintf("block_%d_fileBytes", i)
				raw, exists := c.Get(key)
				if !exists {
					return "", fmt.Errorf("missing file data for video block %d", i)
				}
				fileBytes, ok := raw.([]byte)
				if !ok {
					return "", fmt.Errorf("invalid file data format for video block %d", i)
				}

				// Upload the video
				s3Key, err := s.s3Uploader.UploadFile(
					constant.S3FolderVideo,
					post.Blocks[i].Video.Filename,
					"video/mp4",
					"someUserID",
					fileBytes,
				)
				if err != nil {
					return "", fmt.Errorf("failed to upload video block %d: %w", i, err)
				}

				post.Blocks[i].Video.Filename = s3Key
			}

		case entity.BlockTypeText:
		case entity.BlockTypeHeading:
			// No special action: the block.Text is already filled out
		}
	}

	current_time := time.Now()
	post.CreatedAt = current_time
	post.UpdatedAt = current_time

	post.Status = constant.BlogStatusActive

	// Finally store the post in the DB
	insertedID, err := s.repo.Add(c.Request.Context(), post)
	if err != nil {
		return "", fmt.Errorf("failed to insert blog post: %w", err)
	}

	return insertedID, nil
}

// FindPostsWithRawQuery: parse the raw query params in the service, then build QueryOptions.
func (s *blogPostService) FindPostsWithRawQuery(
	ctx context.Context,
	rawFilters, rawSorts []string,
	rawFields string,
	page, pageSize int,
) ([]entity.BlogPost, error) {
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

	parsedFields := query.ParseFields(rawFields)

	opts := query.QueryOptions{
		Filters: parsedFilters,
		Sorts:   parsedSorts,
		Fields:  parsedFields,
		Skip:    int64((page - 1) * pageSize),
		Limit:   int64(pageSize),
	}

	posts, err := s.repo.FindByQuery(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}

	for pIdx := range posts {
		for bIdx := range posts[pIdx].Blocks {
			block := &posts[pIdx].Blocks[bIdx]
			switch block.Type {
			case entity.BlockTypeImage:
				if block.Image != nil && block.Image.Filename != "" {
					link, presignErr := s.s3Uploader.GeneratePresignedURL(block.Image.Filename, 15*time.Minute)
					if presignErr == nil {
						block.Image.Link = &link
					} else {
						// Handle or log presignErr as needed
					}
				}
			case entity.BlockTypeVideo:
				if block.Video != nil && block.Video.Filename != "" {
					link, presignErr := s.s3Uploader.GeneratePresignedURL(block.Video.Filename, 15*time.Minute)
					if presignErr == nil {
						block.Video.Link = &link
					} else {
						// Handle or log presignErr as needed
					}
				}
			}
		}
	}

	return posts, nil
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

func (s *blogPostService) LoadAllPosts(ctx context.Context) ([]entity.BlogPost, error) {
	posts, err := s.repo.LoadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load all posts: %w", err)
	}
	return posts, nil
}

func (s *blogPostService) SoftDeletePostByRawFilter(
	ctx context.Context, rawFilters []string) error {

	parsedFilters, err := query.ParseFilters(rawFilters)
	if err != nil {
		return fmt.Errorf("failed to parse filters: %w", err)
	}

	filterDoc, _ := query.BuildMongoQuery(query.QueryOptions{
		Filters: parsedFilters,
	})

	fields := bson.M{
		"status":     constant.BlogStatusDeleted,
		"updated_at": time.Now(),
	}

	return s.repo.UpdateFieldsByQuery(ctx, filterDoc, fields)
}
