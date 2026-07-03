package middleware

import (
	"net/http"
	"pendekin_go/config"

	"github.com/gin-gonic/gin"
)

type CORSManager struct {
	cfg config.AppConfig
}

func NewCORSManager(cfg config.AppConfig) *CORSManager {
	return &CORSManager{
		cfg: cfg,
	}
}

func (m *CORSManager) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", m.cfg.FrontendURL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
