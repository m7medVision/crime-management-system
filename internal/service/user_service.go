package service

import (
	"errors"

	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *model.User) (*model.User, error) {
	// Check if username or email already exists
	existingUser, _ := s.userRepo.GetByUsername(user.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, _ = s.userRepo.GetByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *model.User) (*model.User, error) {
	// Check if email is being changed and if it conflicts with another user
	if existingUser, _ := s.userRepo.GetByEmail(user.Email); existingUser != nil && existingUser.ID != user.ID {
		return nil, errors.New("email already in use by another user")
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) ListUsers(offset, limit int) ([]model.User, int64, error) {
	return s.userRepo.List(offset, limit)
}
