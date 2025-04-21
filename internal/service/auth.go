package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"pvz_service/internal/domain"
	"pvz_service/pkg/jwt"
)

// Объявляем ошибки сервиса
var (
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidRole     = errors.New("invalid user role")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, email, password string, role domain.UserRole) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(tokenString string) (*jwt.Claims, error)
}

type authService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthService(repo domain.UserRepository, jwtSecret string) *authService {
	return &authService{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, email, password string, role domain.UserRole) (*domain.User, error) {
	if !role.IsValid() {
		return nil, ErrInvalidRole
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return createdUser, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidPassword
	}

	token, err := jwt.GenerateToken(string(user.Role), s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return jwt.ValidateToken(tokenString, s.jwtSecret)
}