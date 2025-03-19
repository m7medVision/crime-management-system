package service

import (
	"errors"
	"time"

	"github.com/m7medVision/crime-management-system/internal/auth"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type AuthService struct {
	userRepo       *repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", nil, errors.New("account is inactive")
	}

	token, err := auth.GenerateToken(user, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return "", nil, err
	}

	// Update last login time
	now := time.Now()
	user.LastLogin = &now
	err = s.userRepo.Update(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
