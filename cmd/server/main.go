package main

import (
	"os"

	"github.com/capigiba/capiary/internal/config"
	handler "github.com/capigiba/capiary/internal/handler/rest/v1"
	"github.com/capigiba/capiary/internal/infra/db/mongodb"
	"github.com/capigiba/capiary/internal/infra/db/postgres"
	"github.com/capigiba/capiary/internal/infra/storage"
	"github.com/capigiba/capiary/internal/middleware"
	"github.com/capigiba/capiary/internal/repositories"
	"github.com/capigiba/capiary/internal/router"
	"github.com/capigiba/capiary/internal/services"
	"github.com/capigiba/capiary/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	appLogger := logger.NewLogger("Initialize")

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Errorf("config loading error: %w", err)
		os.Exit(1)
	}

	dbPostgresConn, err := postgres.NewPostgresDB(cfg.Database.RdsPostgresURL)
	if err != nil {
		appLogger.Errorf("database initialization error: %w", err)
		os.Exit(1)
	}

	dbMongoConn := mongodb.NewMongoDBClient(cfg.Database.MongodbURI)

	storageClient, err := storage.NewS3Uploader(
		cfg.Storage.AwsAccessKeyID,
		cfg.Storage.AwsSecretKey,
		cfg.Storage.AwsRegion,
		cfg.Storage.AwsBucket,
	)
	if err != nil {
		appLogger.Errorf("Failed to initialize AWS S3 client: %v", err)
	}

	userRepo := repositories.NewUserRepo(dbPostgresConn)
	authUserMiddleware := middleware.NewAuthUserMiddleware(userRepo, cfg.Server.JWTSecret)
	userService := services.NewUserService(userRepo, authUserMiddleware)
	userHandler := handler.NewUserHandler(userService)

	blogRepo := repositories.NewBlogPostRepository(dbMongoConn)
	blogService := services.NewBlogPostService(blogRepo, storageClient)
	blogHandler := handler.NewBlogPostHandler(blogService)

	categoryRepo := repositories.NewCategoryRepository(dbMongoConn)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	swaggerRouter := router.NewSwaggerRouter()

	appRouter := router.NewAppRouter(
		userHandler,
		blogHandler,
		categoryHandler,
		authUserMiddleware,
		swaggerRouter,
	)

	router := gin.Default()

	// Apply CORS middleware with configured settings
	corsConfig := cfg.CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowedOrigins,
		AllowMethods:     corsConfig.AllowedMethods,
		AllowHeaders:     corsConfig.AllowedHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	}))

	apiGroup := router.Group("/api")
	registerAPIRoutes(apiGroup, appRouter)
	registerSwaggerRoutes(router, appRouter)

	port := cfg.Server.Port
	if err := router.Run(":" + port); err != nil {
		appLogger.Errorf("Failed to start the server: %v", err)
		os.Exit(1)
	}
}

func registerAPIRoutes(group *gin.RouterGroup, appRouter *router.AppRouter) {
	appRouter.RegisterUserRoutes(group)
	appRouter.RegisterBlogRoutes(group)
	appRouter.RegisterCategoryRoutes(group)
}

func registerSwaggerRoutes(router *gin.Engine, appRouter *router.AppRouter) {
	swaggerGroup := router.Group("/")
	appRouter.RegisterSwaggerRoutes(swaggerGroup)
}
