package main

import (
	"log"
	"pendekin_go/config"
	"pendekin_go/internal/database"
	"pendekin_go/internal/handler"
	"pendekin_go/internal/middleware"
	"pendekin_go/internal/repository"
	"pendekin_go/internal/usecase"
	"pendekin_go/pkg/jwt"
	"pendekin_go/pkg/logger"
	"pendekin_go/pkg/validation"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Load JWT Manager
	JWT := jwt.NewJWT(&cfg.App)

	// Init Logger
	logger.Init(cfg.App.AppEnv)

	// Register Validators
	validation.RegisterValidators()

	// Init Middleware
	authMiddleware := middleware.NewAuthMiddleware(JWT)
	corsMiddleware := middleware.NewCORSManager(cfg.App)

	// Init Repository
	shortUrlRepo := repository.NewShortUrlRepository(db.DB)
	clickLogRepo := repository.NewClickLogRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Init Usecase
	shortUrlUseCase := usecase.NewShortUrlUsecase(shortUrlRepo, clickLogRepo, &cfg.App)
	userUseCase := usecase.NewUserUseCase(userRepo, JWT)

	// Init Handlers
	healthHandler := handler.NewHealthHandler(db)
	shortUrlHandler := handler.NewShortUrlHandler(shortUrlUseCase)
	userHandler := handler.NewUserHandler(userUseCase, &cfg.App)

	// Setup Router
	router := gin.Default()
	router.Use(corsMiddleware.CORS())

	api := router.Group("/api")

	// Public
	api.GET("/healthcheck", healthHandler.HealthCheck)
	router.GET("/:alias", shortUrlHandler.Redirect)

	// Auth
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)

	// Short Url
	api.Use(authMiddleware.Handle())
	{
		api.POST("/short-urls/create", shortUrlHandler.Create)
		api.DELETE("/short-urls/:id", shortUrlHandler.Delete)

		api.GET("/auth/me", userHandler.FetchMe)
		api.POST("/auth/logout", userHandler.Logout)
	}

	router.Run(":8080")
}
