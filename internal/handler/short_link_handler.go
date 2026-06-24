package handler

import (
	"errors"
	"net/http"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": validation.FormatValidationErrors(ve),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

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
