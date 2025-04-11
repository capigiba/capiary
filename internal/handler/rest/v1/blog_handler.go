package handler

import (
	"net/http"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/services"
	"github.com/gin-gonic/gin"
)

type BlogPostHandler struct {
	service services.BlogPostService
}

func NewBlogPostHandler(service services.BlogPostService) *BlogPostHandler {
	return &BlogPostHandler{service: service}
}

// Create a new blog post
func (h *BlogPostHandler) CreateBlogPostHandler(c *gin.Context) {
	var input entity.BlogPost
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	insertedID, err := h.service.CreatePost(c.Request.Context(), input)
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
