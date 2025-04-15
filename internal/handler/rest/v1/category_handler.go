package handler

import (
	"net/http"

	"github.com/capigiba/capiary/internal/domain/entity"
	"github.com/capigiba/capiary/internal/services"
	"github.com/gin-gonic/gin"
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
	var category entity.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
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

	categories, err := h.service.Find(c.Request.Context(), rawFilters, rawSorts, rawFields)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategoryHandler updates a category using raw filters in query params and JSON body as the update data.
// e.g. PATCH /categories?filter=id__==__<someObjectId>
func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	rawFilters := c.QueryArray("filter")

	var update entity.Category
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if err := h.service.UpdateByRawFilter(c.Request.Context(), rawFilters, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category: " + err.Error()})
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
