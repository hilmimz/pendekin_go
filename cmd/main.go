package main

import (
	"log"
	"net/http"
	"pendekin_go/config"
	"pendekin_go/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		status := db.HealthCheck()
		if status != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Database connection failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
			"db":      "Connected",
		})
	})

	router.Run(":8080")
}
