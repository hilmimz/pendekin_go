package main

import (
	"log"
	"pendekin_go/config"
	"pendekin_go/internal/database"
	"pendekin_go/internal/handler"
	"pendekin_go/internal/repository"
	"pendekin_go/internal/usecase"
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

	// Register Validators
	validation.RegisterValidators()

	// Init Repository
	shortLinkRepo := repository.NewShortLinkRepository(db.DB)
	clickLogRepo := repository.NewClickLogRepository(db.DB)

	// Init Usecase
	shortLinkUseCase := usecase.NewShortLinkUseCase(shortLinkRepo, clickLogRepo, &cfg.App)

	// Init Handlers
	healthHandler := handler.NewHealthHandler(db)
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkUseCase)

	// Setup Router
	router := gin.Default()
	api := router.Group("/api")
	api.GET("/healthcheck", healthHandler.HealthCheck)
	api.POST("/short-links/create", shortLinkHandler.Create)
	router.GET("/:alias", shortLinkHandler.Redirect)

	// Short Link Route

	router.Run(":8080")
}
