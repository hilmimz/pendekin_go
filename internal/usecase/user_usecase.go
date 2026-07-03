package usecase

import (
	"errors"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/errs"
	"pendekin_go/pkg/hash"
	"pendekin_go/pkg/logger"

	"go.uber.org/zap"
)

type UserUseCase struct {
	UserRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{UserRepo: userRepo}
}

func (uc *UserUseCase) Register(req *domain.UserRegisterRequest) (*domain.UserRegisterResponse, *errs.Error) {
	user, err := uc.UserRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, domain.ErrEmailNotFound) {
		logger.Log.Error("failed to find user by email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errs.Internal("failed to find user by email", err)
	}
	if user != nil {
		logger.Log.Error("email already exists",
			zap.String("email", req.Email),
		)
		return nil, errs.Conflict("email already exists", nil)
	}

	// password hashing
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error("failed to hash password",
			zap.Error(err),
		)
		return nil, errs.Internal("failed to hash password", err)
	}
	err = uc.UserRepo.Create(&domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		logger.Log.Error("failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errs.Internal("failed to create user", err)
	}
	registeredUser, err := uc.UserRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, domain.ErrEmailNotFound) {
		logger.Log.Error("failed to find user by email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errs.Internal("failed to find user by email", err)
	}
	res := domain.UserRegisterResponse{
		ID:        registeredUser.ID,
		Name:      registeredUser.Name,
		Email:     registeredUser.Email,
		CreatedAt: registeredUser.CreatedAt,
	}
	logger.Log.Info("user registered successfully",
		zap.Int("id", res.ID),
		zap.String("name", res.Name),
		zap.String("email", res.Email),
	)
	return &res, nil
}
