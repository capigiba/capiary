package router

import (
	handler "github.com/capigiba/capiary/internal/handler/rest/v1"
	"github.com/capigiba/capiary/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	userController     *handler.UserHandler
	blogController     *handler.BlogPostHandler
	categoryController *handler.CategoryHandler
	authMiddleware     *middleware.AuthUserMiddleware
	swaggerRouter      *SwaggerRouter
}

func NewAppRouter(
	userController *handler.UserHandler,
	blogController *handler.BlogPostHandler,
	categoryController *handler.CategoryHandler,
	authMiddleware *middleware.AuthUserMiddleware,
	swaggerRouter *SwaggerRouter) *AppRouter {
	return &AppRouter{
		userController:     userController,
		blogController:     blogController,
		categoryController: categoryController,
		authMiddleware:     authMiddleware,
		swaggerRouter:      swaggerRouter,
	}
}

// RegisterUserRoutes sets up the routes for user-related operations
func (a *AppRouter) RegisterUserRoutes(r *gin.RouterGroup) {
	public := r.Group("/users")
	{
		public.POST("/register", a.userController.RegisterUser)
		public.POST("/login", a.userController.Login)
	}

	protected := r.Group("/users")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.PUT("/:user_id/change-password", a.userController.ChangePassword)
	}
}

func (a *AppRouter) RegisterBlogRoutes(r *gin.RouterGroup) {
	protected := r.Group("/blog")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/posts", a.blogController.CreateBlogPostHandler)
		protected.GET("/posts", a.blogController.FindBlogPostsHandler)
		protected.PUT("/posts", a.blogController.UpdateBlogPostHandler)
		protected.GET("/posts/all", a.blogController.LoadAllPostsHandler)
		protected.DELETE("/posts", a.blogController.DeleteBlogPostHandler)
	}
}

func (a *AppRouter) RegisterCategoryRoutes(r *gin.RouterGroup) {
	public := r.Group("/categories")
	{
		public.POST("/all", a.categoryController.LoadAllCategoriesHandler)
	}

	protected := r.Group("/categories")
	protected.Use(a.authMiddleware.MustAuth())
	{
		protected.POST("/", a.categoryController.CreateCategoryHandler)
		protected.GET("/", a.categoryController.FindCategoriesHandler)
		protected.PUT("/", a.categoryController.UpdateCategoryHandler)
	}
}

// RegisterSwaggerRoutes sets up the route for Swagger API documentation
func (a *AppRouter) RegisterSwaggerRoutes(r *gin.RouterGroup) {
	// Check if SwaggerRouter is initialized before registering
	if a.swaggerRouter != nil {
		a.swaggerRouter.Register(r)
	}
}
