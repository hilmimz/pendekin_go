package handler

import (
	"net/http"
	"pendekin_go/internal/database"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	db *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := h.db.HealthCheck()
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
}
