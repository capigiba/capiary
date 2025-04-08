package router

import (
	"fmt"

	"github.com/capigiba/capiary/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerConfig struct describes configure for the Swagger API endpoint
type SwaggerConfig struct {
	Show     bool
	Protocol string
	Host     string
	Address  string
}

// NewSwaggerRouter initializes and returns a new SwaggerRouter
func NewSwaggerRouter() *SwaggerRouter {
	// Set the configuration values directly in the file
	config := &SwaggerConfig{
		Show:     true,
		Protocol: "http",
		Host:     "localhost",
		Address:  ":8080",
	}

	return &SwaggerRouter{
		config: config,
	}
}

// SwaggerRouter swagger api router
type SwaggerRouter struct {
	config *SwaggerConfig
}

// Register register swagger api router
func (a *SwaggerRouter) Register(r *gin.RouterGroup) {
	if a.config.Show {
		a.InitSwaggerDocs()
		gofmt := fmt.Sprintf("%s://%s%s/swagger/doc.json", a.config.Protocol, a.config.Host, a.config.Address)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(gofmt)))
	}
}

// InitSwaggerDocs init swagger docs
func (a *SwaggerRouter) InitSwaggerDocs() {
	docs.SwaggerInfo.Title = "Borrow_books API"
	docs.SwaggerInfo.Description = "API documentation for Borrow_books"
	docs.SwaggerInfo.Version = "v0.0.1"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s%s", a.config.Host, a.config.Address)
	docs.SwaggerInfo.BasePath = "/api"
}
