package usecase

import (
	"errors"
	"pendekin_go/internal/domain"
	"pendekin_go/pkg/errs"
	"pendekin_go/pkg/hash"
	"pendekin_go/pkg/jwt"
	"pendekin_go/pkg/logger"

	"go.uber.org/zap"
)

type UserUseCase struct {
	UserRepo domain.UserRepository
	JWT      *jwt.JWT
}

func NewUserUseCase(userRepo domain.UserRepository, jwt *jwt.JWT) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
		JWT:      jwt,
	}
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

	token, err := uc.JWT.GenerateToken(registeredUser.ID, registeredUser.Email, registeredUser.Name)
	if err != nil {
		logger.Log.Error("failed to generate token",
			zap.Error(err),
		)
		return nil, errs.Internal("failed to generate token", err)
	}

	res := domain.UserRegisterResponse{
		ID:        registeredUser.ID,
		Name:      registeredUser.Name,
		Email:     registeredUser.Email,
		CreatedAt: registeredUser.CreatedAt,
		Token:     token,
	}

	logger.Log.Info("user registered successfully",
		zap.Int("id", res.ID),
		zap.String("name", res.Name),
		zap.String("email", res.Email),
	)
	return &res, nil
}

func (uc *UserUseCase) Login(req *domain.UserLoginRequest) (*domain.UserLoginResponse, *errs.Error) {
	user, err := uc.UserRepo.FindByEmail(req.Email)
	if err != nil && !errors.Is(err, domain.ErrEmailNotFound) {
		logger.Log.Error("failed to find user by email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, errs.Internal("failed to find user by email", err)
	}
	if user == nil {
		return nil, errs.Unauthorized("invalid email or password", nil)
	}

	pw := hash.CheckPassword(user.Password, req.Password)
	if !pw {
		return nil, errs.Unauthorized("invalid email or password", nil)
	}

	token, err := uc.JWT.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		logger.Log.Error("failed to generate token",
			zap.Error(err),
		)
		return nil, errs.Internal("failed to generate token", err)
	}

	res := &domain.UserLoginResponse{
		Email: user.Email,
		Name:  user.Name,
		Token: token,
	}

	return res, nil
}

func (uc *UserUseCase) FetchMe(req *domain.FetchMeRequest) (*domain.FetchMeResponse, *errs.Error) {
	user, err := uc.UserRepo.FindByID(req.UserID)
	if err != nil && !errors.Is(err, domain.ErrEmailNotFound) {
		logger.Log.Error("failed to find user by id",
			zap.Int("Id", req.UserID),
			zap.Error(err),
		)
		return nil, errs.Internal("failed to find user by id", err)
	}
	if user == nil {
		return nil, errs.Unauthorized("invalid id", nil)
	}

	res := &domain.FetchMeResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return res, nil
}
