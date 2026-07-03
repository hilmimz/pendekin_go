package handler

import (
	"errors"
	"net/http"
	"pendekin_go/internal/domain"
	"pendekin_go/internal/usecase"
	"pendekin_go/pkg/response"
	"pendekin_go/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
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
	response.ResponseOK(c, http.StatusOK, "login successful", resp)
}
