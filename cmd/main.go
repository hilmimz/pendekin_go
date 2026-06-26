package main

import (
	"log"
	"pendekin_go/config"
	"pendekin_go/internal/database"
	"pendekin_go/internal/handler"
	"pendekin_go/internal/repository"
	"pendekin_go/internal/usecase"
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

	// Init Logger
	logger.Init(cfg.App.AppEnv)

	// Register Validators
	validation.RegisterValidators()

	// Init Repository
	shortUrlRepo := repository.NewShortUrlRepository(db.DB)
	clickLogRepo := repository.NewClickLogRepository(db.DB)

	// Init Usecase
	shortUrlUseCase := usecase.NewShortUrlUsecase(shortUrlRepo, clickLogRepo, &cfg.App)

	// Init Handlers
	healthHandler := handler.NewHealthHandler(db)
	shortUrlHandler := handler.NewShortUrlHandler(shortUrlUseCase)

	// Setup Router
	router := gin.Default()
	api := router.Group("/api")
	api.GET("/healthcheck", healthHandler.HealthCheck)
	api.POST("/short-urls/create", shortUrlHandler.Create)
	api.DELETE("/short-urls/:id", shortUrlHandler.Delete)

	router.GET("/:alias", shortUrlHandler.Redirect)

	// Short Url Route

	router.Run(":8080")
}
