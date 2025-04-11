package handler

import (
	"net/http"

	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/domain/request"
	"github.com/capigiba/capiary/internal/services"
	"github.com/gin-gonic/gin"
)

type BlogPostHandler struct {
	service services.BlogPostService
}

func NewBlogPostHandler(service services.BlogPostService) *BlogPostHandler {
	return &BlogPostHandler{service: service}
}

func (h *BlogPostHandler) CreateBlogPostHandler(c *gin.Context) {
	// parse incoming JSON into our request DTO
	var req request.CreateBlogPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	post := entity.BlogPost{
		Title: req.Title,
	}

	// DTO: CreateBlockRequest -> entity.Block
	var blocks []entity.Block
	for i, blockReq := range req.Blocks {
		block := entity.Block{
			ID:    i + 1,
			Order: blockReq.Order,
		}

		switch blockReq.Type {
		case constant.MediaTypeText:
			block.Type = entity.BlockTypeText
			// convert paragraphs
			textBlock := entity.TextBlock{}
			for _, paraReq := range blockReq.Paragraphs {
				paragraph := entity.Paragraph{
					ID:   paraReq.ID,
					Text: paraReq.Text,
				}

				// convert formats
				for _, formatReq := range paraReq.Formats {
					paragraph.Formats = append(paragraph.Formats, entity.Format{
						Type:      entity.FormatType(formatReq.Type), // cast the string
						Start:     formatReq.Start,
						End:       formatReq.End,
						Hyperlink: formatReq.Hyperlink,
					})
				}

				textBlock.Paragraphs = append(textBlock.Paragraphs, paragraph)
			}
			block.Text = &textBlock

		case constant.MediaTypeImage:
			block.Type = entity.BlockTypeImage
			block.Image = &entity.ImageBlock{
				Filename: blockReq.Filename,
			}

		case constant.MediaTypeVideo:
			block.Type = entity.BlockTypeVideo
			block.Video = &entity.VideoBlock{
				Filename: blockReq.Filename,
			}
		}

		blocks = append(blocks, block)
	}

	post.Blocks = blocks

	insertedID, err := h.service.CreatePostWithFiles(c, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
}

// Find blog posts with raw filter/sort/fields
func (h *BlogPostHandler) FindBlogPostsHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter") // e.g. ["age__gt__30", "title__==__Hello"]
	rawSorts := c.QueryArray("sort")     // e.g. ["age__desc", "title__asc"]
	rawFields := c.Query("fields")       // e.g. "id,title"

	posts, err := h.service.FindPostsWithRawQuery(c.Request.Context(), rawFilters, rawSorts, rawFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// Update a blog post using a filter
func (h *BlogPostHandler) UpdateBlogPostHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter") // e.g. ["_id__==__<ObjectID>"]

	var update entity.BlogPost
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	err := h.service.UpdatePostByRawFilter(c.Request.Context(), rawFilters, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "post updated"})
}
