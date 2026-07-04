package handler

import (
	"errors"
	"net/http"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/logger"
	"pendekin_go/pkg/response"
	"pendekin_go/pkg/validation"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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
			logger.Log.Error("validation error", zap.Error(err))
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}
		logger.Log.Error("invalid request body", zap.Error(err))
		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Log.Error("unauthorized: user_id not found in context")
		response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userIdInt := userId.(int)
	req.UserID = &userIdInt

	resp, err := s.shortUrlUseCase.CreateShortUrl(&req)
	if err != nil {
		logger.Log.Error("failed to create short url", zap.Error(err), zap.Any("request", req))
		response.ResponseNOK(c, err.Code, err.Message, nil)
		return
	}
	logger.Log.Info("short url created successfully", zap.Int("user_id", userIdInt), zap.Int("short_url_id", resp.ID))
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
		logger.Log.Error("failed to redirect short url", zap.Error(err), zap.String("alias", alias))
		// Should redirect to frontend 404 or return html
		response.ResponseNOK(c, err.Code, err.Message, nil)
		return
	}
	logger.Log.Info("short url redirected successfully", zap.String("alias", alias), zap.String("original_url", resp.OriginalURL))
	c.Redirect(http.StatusMovedPermanently, resp.OriginalURL)
}

func (s *ShortUrlHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("invalid id parameter", zap.Error(err), zap.String("id_param", c.Param("id")))
		response.ResponseNOK(c, http.StatusBadRequest, "invalid id", nil)
		return
	}

	userId, ok := c.Get("user_id")
	if !ok {
		logger.Log.Error("unauthorized: user_id not found in context")
		response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userIdInt := userId.(int)
	pointerId := &userIdInt

	req := domain.DeleteShortUrlRequest{
		ID:     id,
		UserID: pointerId,
	}

	resp, errs := s.shortUrlUseCase.DeleteShortUrl(&req)
	if errs != nil {
		logger.Log.Error("failed to delete short url", zap.Error(errs), zap.Int("short_url_id", id), zap.Int("user_id", userIdInt))
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	logger.Log.Info("short url deleted successfully", zap.Int("short_url_id", id), zap.Int("user_id", userIdInt))
	response.ResponseOK(c, 200, "Short url deleted successfully", resp)
}
