package service

import (
	"context"
	"errors"
	"time"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/middleware"
	"github.com/codewithwan/gopilot/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
}

type authService struct {
	userRepo       repository.UserRepository
	jwtMiddleware  *middleware.JWTMiddleware
	jwtExpiration  time.Duration
	logger         *zap.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtMiddleware *middleware.JWTMiddleware, jwtExpiration time.Duration, logger *zap.Logger) AuthService {
	return &authService{
		userRepo:      userRepo,
		jwtMiddleware: jwtMiddleware,
		jwtExpiration: jwtExpiration,
		logger:        logger,
	}
}

func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err), zap.String("username", req.Username))
		return nil, err
	}

	s.logger.Info("user registered", zap.Int64("user_id", created.ID), zap.String("username", created.Username))
	return created, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err), zap.String("username", req.Username))
		return nil, ErrInvalidCredentials
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warn("invalid password attempt", zap.String("username", req.Username))
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwtMiddleware.GenerateToken(user.ID, user.Username, s.jwtExpiration)
	if err != nil {
		s.logger.Error("failed to generate token", zap.Error(err), zap.Int64("user_id", user.ID))
		return nil, err
	}

	s.logger.Info("user logged in", zap.Int64("user_id", user.ID), zap.String("username", user.Username))
	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
