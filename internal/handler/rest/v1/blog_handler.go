package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/capigiba/capiary/internal/domain/constant"
	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/domain/request"
	"github.com/capigiba/capiary/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogPostHandler struct {
	service services.BlogPostService
}

func NewBlogPostHandler(service services.BlogPostService) *BlogPostHandler {
	return &BlogPostHandler{service: service}
}

func (h *BlogPostHandler) CreateBlogPostHandler(c *gin.Context) {

	// 1) Parse the "metadata" field from the multipart form.
	metadataJSON := c.PostForm("metadata")
	if metadataJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing metadata in form"})
		return
	}

	// 2) Unmarshal metadata into your request struct.
	var req request.CreateBlogPostRequest
	if err := json.Unmarshal([]byte(metadataJSON), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metadata JSON: " + err.Error()})
		return
	}

	// 3) Now handle each block’s file, if any.
	//    For block i, the file field on the front-end is "block_i_file"
	for i, blockReq := range req.Blocks {
		if blockReq.Type == constant.MediaTypeImage || blockReq.Type == constant.MediaTypeVideo {
			fieldName := fmt.Sprintf("block_%d_file", i)
			fileHeader, err := c.FormFile(fieldName)
			if err != nil {
				// If we get an error here, it might just mean the user didn’t attach a file
				// Or it could be a real problem. Decide how to handle that:
				if err == http.ErrMissingFile {
					// Possibly no file for that block
					continue
				}
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// Read the file bytes
			f, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer f.Close()

			fileBytes, err := io.ReadAll(f)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 4) Stash these bytes in the Gin context for your service code to pick up
			c.Set(fmt.Sprintf("block_%d_fileBytes", i), fileBytes)
		}
	}

	// 5) Now that all data is in place, do the normal create:
	//    This calls h.service.CreatePostWithFiles(c, post) in your code
	post := entity.BlogPost{
		ID:    primitive.NewObjectID(),
		Title: req.Title,
		// AuthorID:   req.AuthorID,
		// Categories: req.Categories,
	}

	// Rebuild the blocks array from req.Blocks, just like your snippet does:
	var blocks []entity.Block
	for i, blockReq := range req.Blocks {
		id := blockReq.ID
		if id == 0 {
			id = i + 1
		}
		block := entity.Block{
			ID:    id,
			Order: blockReq.Order,
		}
		switch blockReq.Type {
		case constant.MediaTypeText:
			block.Type = entity.BlockTypeText
			textBlock := entity.TextBlock{Paragraphs: []entity.Paragraph{}}
			for j, paraReq := range blockReq.Paragraphs {
				para := entity.Paragraph{
					ID:    j + 1,
					Text:  paraReq.Text,
					Align: paraReq.Align,
				}
				for _, f := range paraReq.Formats {
					para.Formats = append(para.Formats, entity.Format{
						Type:      entity.FormatType(f.Type),
						Start:     f.Start,
						End:       f.End,
						Hyperlink: f.Hyperlink,
					})
				}
				textBlock.Paragraphs = append(textBlock.Paragraphs, para)
			}
			block.Text = &textBlock
		case constant.MediaTypeHeading:
			if blockReq.HeadingLevel == nil || blockReq.Text == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "headingLevel and text are required for heading blocks",
				})
				return
			}
			block.Type = entity.BlockTypeHeading
			block.Heading = &entity.HeadingBlock{
				Level: *blockReq.HeadingLevel,
				Text:  *blockReq.Text,
			}
		case constant.MediaTypeImage:
			block.Type = entity.BlockTypeImage
			block.Image = &entity.ImageBlock{
				Filename: blockReq.Filename, // from JSON
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

	// Actually call your service
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

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	posts, err := h.service.FindPostsWithRawQuery(c.Request.Context(), rawFilters, rawSorts, rawFields, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": posts,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"count":     len(posts),
		},
	})
}

// Update a blog post using a filter
func (h *BlogPostHandler) UpdateBlogPostHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter") // e.g. ["_id__==__<ObjectID>"]

	if len(rawFilters) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing filter"})
		return
	}

	metaJSON := c.PostForm("metadata")
	if metaJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing metadata field"})
		return
	}

	var req request.CreateBlogPostRequest // reuse same shape
	if err := json.Unmarshal([]byte(metaJSON), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metadata: " + err.Error()})
		return
	}

	for i, b := range req.Blocks {
		if b.Type != constant.MediaTypeImage && b.Type != constant.MediaTypeVideo {
			continue
		}
		if fh, err := c.FormFile(fmt.Sprintf("block_%d_file", i)); err == nil {
			f, _ := fh.Open()
			defer f.Close()
			bytes, _ := io.ReadAll(f)
			c.Set(fmt.Sprintf("block_%d_fileBytes", i), bytes)
			// keep original filename for MIME check / ext if needed
			c.Set(fmt.Sprintf("block_%d_origFilename", i), fh.Filename)
		}
	}

	post := entity.BlogPost{
		Title: req.Title,
	}
	var blocks []entity.Block
	for i, b := range req.Blocks {
		id := b.ID
		if id == 0 {
			id = i + 1
		}
		blk := entity.Block{ID: id, Order: b.Order}
		switch b.Type {
		case constant.MediaTypeText:
			blk.Type = entity.BlockTypeText
			var paras []entity.Paragraph
			for j, p := range b.Paragraphs {
				para := entity.Paragraph{
					ID:    j + 1,
					Text:  p.Text,
					Align: p.Align,
				}
				for _, f := range p.Formats {
					para.Formats = append(para.Formats, entity.Format{
						Type:      entity.FormatType(f.Type),
						Start:     f.Start,
						End:       f.End,
						Hyperlink: f.Hyperlink,
					})
				}
				paras = append(paras, para)
			}
			blk.Text = &entity.TextBlock{Paragraphs: paras}

		case constant.MediaTypeHeading:
			if b.HeadingLevel == nil || b.Text == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "heading requires level & text"})
				return
			}
			blk.Type = entity.BlockTypeHeading
			blk.Heading = &entity.HeadingBlock{Level: *b.HeadingLevel, Text: *b.Text}

		case constant.MediaTypeImage:
			blk.Type = entity.BlockTypeImage
			blk.Image = &entity.ImageBlock{Filename: b.Filename}

		case constant.MediaTypeVideo:
			blk.Type = entity.BlockTypeVideo
			blk.Video = &entity.VideoBlock{Filename: b.Filename}
		}
		blocks = append(blocks, blk)
	}
	post.Blocks = blocks

	err := h.service.UpdatePostByRawFilter(c, rawFilters, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "post updated"})
}

func (h *BlogPostHandler) LoadAllPostsHandler(c *gin.Context) {
	posts, err := h.service.LoadAllPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *BlogPostHandler) DeleteBlogPostHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter") // same format as update

	if len(rawFilters) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing filter"})
		return
	}

	if err := h.service.SoftDeletePostByRawFilter(c.Request.Context(), rawFilters); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "post deleted (soft)"})
}
