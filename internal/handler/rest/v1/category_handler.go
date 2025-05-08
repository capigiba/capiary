package handler

import (
	"net/http"
	"strconv"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/domain/request"
	"github.com/capigiba/capiary/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryHandler wraps the CategoryService for use in HTTP handlers.
type CategoryHandler struct {
	service services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler with the given CategoryService.
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategoryHandler creates a new category using JSON request data.
func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	category := entity.Category{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
	}

	insertedID, err := h.service.Create(c, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
}

// FindCategoriesHandler finds categories based on raw filter/sort/fields in query params.
// e.g. GET /categories?filter=name__==__someName&sort=name__asc&fields=name,description
func (h *CategoryHandler) FindCategoriesHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter")
	rawSorts := c.QueryArray("sort")
	rawFields := c.Query("fields")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	categories, err := h.service.Find(c.Request.Context(), rawFilters, rawSorts, rawFields, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
		"meta": gin.H{
			"page":      page,
			"page_size": pageSize,
			"count":     len(categories),
		},
	})
}

// UpdateCategoryHandler updates a category using raw filters in query params and JSON body as the update data.
// e.g. PATCH /categories?filter=id__==__<someObjectId>
func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter")
	if len(rawFilters) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing filter"})
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	update := entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.service.UpdateByRawFilter(c, rawFilters, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "category updated"})
}

// LoadAllCategoriesHandler loads all categories in ascending order by name.
func (h *CategoryHandler) LoadAllCategoriesHandler(c *gin.Context) {
	categories, err := h.service.LoadAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load categories: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}
