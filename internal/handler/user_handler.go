package handler

import (
	"errors"
	"net/http"
	"pendekin_go/config"
	"pendekin_go/internal/domain"
	"pendekin_go/internal/usecase"
	"pendekin_go/pkg/response"
	"pendekin_go/pkg/validation"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}

		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, errs := h.userUseCase.Register(&req)
	if errs != nil {
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	response.ResponseOK(c, http.StatusCreated, "user registered successfully", resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			response.ResponseNOK(c, http.StatusBadRequest, "validation error", validation.FormatValidationErrors(ve))
			return
		}

		response.ResponseNOK(c, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	resp, errs := h.userUseCase.Login(&req)
	if errs != nil {
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
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

func (h *UserHandler) FetchMe(c *gin.Context) {
	var req domain.FetchMeRequest
	id, err := c.Cookie("id")
	if err != nil {
		response.ResponseNOK(c, http.StatusUnauthorized, "invalid cookie", nil)
		return
	}
	req.UserID, _ = strconv.Atoi(id)

	resp, errs := h.userUseCase.FetchMe(&req)
	if errs != nil {
		response.ResponseNOK(c, errs.Code, errs.Message, nil)
		return
	}
	response.ResponseOK(c, http.StatusOK, "fetch me successful", resp)
}
