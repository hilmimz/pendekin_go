package handler

import (
	"errors"
	"net/http"
	"pendekin_go/config"
	"pendekin_go/internal/domain"
	"pendekin_go/internal/usecase"
	"pendekin_go/pkg/logger"
	"pendekin_go/pkg/response"
	"pendekin_go/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	cfg         *config.AppConfig
}

func NewUserHandler(userUseCase *usecase.UserUseCase, cfg *config.AppConfig) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		cfg:         cfg,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req domain.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			logger.Log.Error("validation error", zap.Error(err), zap.Any("request", req))
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}

		logger.Log.Error("invalid request body", zap.Error(err), zap.Any("request", req))
		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, errs := h.userUseCase.Register(&req)
	if errs != nil {
		logger.Log.Error("failed to register user", zap.Error(errs), zap.String("email", req.Email))
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	logger.Log.Info("user registered successfully", zap.Int("user_id", resp.ID), zap.String("email", resp.Email))
	c.SetCookie(
		"token",
		resp.Token,
		h.cfg.JWTExpiresIn, // Change if remember me for 14 days is implemented
		"/",
		"",
		true,
		true,
	)
	response.ResponseOK(c, http.StatusCreated, "user registered successfully", gin.H{
		"id":         resp.ID,
		"name":       resp.Name,
		"email":      resp.Email,
		"created_at": resp.CreatedAt,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			logger.Log.Error("validation error", zap.Error(err), zap.Any("request", req))
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}

		logger.Log.Error("invalid request body", zap.Error(err), zap.Any("request", req))
		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, errs := h.userUseCase.Login(&req)
	if errs != nil {
		logger.Log.Error("failed to login", zap.Error(errs), zap.String("email", req.Email))
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	logger.Log.Info("login successful", zap.String("email", resp.Email), zap.String("name", resp.Name))
	c.SetCookie(
		"token",
		resp.Token,
		h.cfg.JWTExpiresIn, // Change if remember me for 14 days is implemented
		"/",
		"",
		true,
		true,
	)
	response.ResponseOK(c, http.StatusOK, "login successful", gin.H{
		"email": resp.Email,
		"name":  resp.Name,
	})
}

func (h *UserHandler) Logout(c *gin.Context) {
	userId, ok := c.Get("user_id")
	if !ok {
		logger.Log.Error("unauthorized: user_id not found in context")
		response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	logger.Log.Info("logout successful", zap.Int("user_id", userId.(int)))
	c.SetCookie("token", "", -1, "/", "", true, true)
	response.ResponseOK(c, http.StatusOK, "logout success", nil)
}

func (h *UserHandler) FetchMe(c *gin.Context) {
	var req domain.FetchMeRequest
	userId, ok := c.Get("user_id")
	if !ok {
		logger.Log.Error("unauthorized: user_id not found in context")
		response.ResponseNOK(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userIdInt := userId.(int)
	req.UserID = userIdInt

	resp, errs := h.userUseCase.FetchMe(&req)
	if errs != nil {
		logger.Log.Error("failed to fetch user data", zap.Error(errs), zap.Int("user_id", userIdInt))
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	logger.Log.Info("fetch me successful", zap.Int("user_id", userIdInt))
	response.ResponseOK(c, http.StatusOK, "fetch me successful", resp)
}
