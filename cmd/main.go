package main

import (
	"log"
	"pendekin_go/config"
	"pendekin_go/internal/database"
	"pendekin_go/internal/handler"

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

	// Init Handlers
	healthHandler := handler.NewHealthHandler(db)

	// Setup Router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/healthcheck", healthHandler.HealthCheck)

	router.Run(":8080")
}
