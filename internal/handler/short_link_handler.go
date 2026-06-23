package handler

import (
	"fmt"
	"pendekin_go/internal/domain"

	"github.com/gin-gonic/gin"
)

type ShortLinkHandler struct {
	shortLinkUseCase domain.ShortLinkUsecase
}

func NewShortLinkHandler(shortLinkUseCase domain.ShortLinkUsecase) *ShortLinkHandler {
	return &ShortLinkHandler{
		shortLinkUseCase: shortLinkUseCase,
	}
}

func (s *ShortLinkHandler) Create(c *gin.Context) {
	var req domain.CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	fmt.Println("Request Body:", req)
	resp, err := s.shortLinkUseCase.CreateShortLink(&req)
	if err != nil {
		c.JSON(err.Code, gin.H{
			"code":    err.Code,
			"message": err.Message,
			"error":   err.Err.Error(),
		})
		return
	}
	c.JSON(201, resp)
}
