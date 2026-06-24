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
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}

		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, err := s.shortLinkUseCase.CreateShortLink(&req)
	if err != nil {
		response.ResponseNOK(c, err.Code, err.Message, nil)
		return
	}
	response.ResponseOK(c, 201, "Short link created successfully", resp)
}

func (s *ShortLinkHandler) Redirect(c *gin.Context) {
	var req domain.RedirectShortLinkRequest
	alias := c.Param("alias")

	req = domain.RedirectShortLinkRequest{
		Alias: &alias,
	}

	resp, err := s.shortLinkUseCase.RedirectShortLink(&req)
	if err != nil {
		// Should redirect to frontend 404 or return html
		response.ResponseNOK(c, err.Code, err.Message, nil)
		return
	}

	c.Redirect(http.StatusMovedPermanently, resp.OriginalURL)
}
