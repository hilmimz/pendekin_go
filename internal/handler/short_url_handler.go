package handler

import (
	"errors"
	"net/http"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/response"
	"pendekin_go/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ShortUrlHandler struct {
	shortUrlUseCase domain.ShortUrlUsecase
}

func NewShortUrlHandler(shortUrlUseCase domain.ShortUrlUsecase) *ShortUrlHandler {
	return &ShortUrlHandler{
		shortUrlUseCase: shortUrlUseCase,
	}
}

func (s *ShortUrlHandler) Create(c *gin.Context) {
	var req domain.CreateShortUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}

		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, err := s.shortUrlUseCase.CreateShortUrl(&req)
	if err != nil {
		response.ResponseNOK(c, err.Code, err.Message, nil)
		return
	}
	response.ResponseOK(c, 201, "Short url created successfully", resp)
}

func (s *ShortUrlHandler) Redirect(c *gin.Context) {
	var req domain.RedirectShortUrlRequest
	alias := c.Param("alias")
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referer := c.GetHeader("Referer")

	req = domain.RedirectShortUrlRequest{
		Alias:     &alias,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Referer:   referer,
	}

	resp, err := s.shortUrlUseCase.RedirectShortUrl(&req)
	if err != nil {
		// Should redirect to frontend 404 or return html
		response.ResponseNOK(c, err.Code, err.Message, nil)
		return
	}

	c.Redirect(http.StatusMovedPermanently, resp.OriginalURL)
}
