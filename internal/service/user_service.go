package service

import (
	"github.com/m7medVision/crime-management-system/internal/auth"
	"github.com/m7medVision/crime-management-system/internal/dto"
	"github.com/m7medVision/crime-management-system/internal/model"
	"github.com/m7medVision/crime-management-system/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(createUserDTO *dto.CreateUserDTO) (*model.User, error) {
	hashedPassword, err := auth.HashPassword(createUserDTO.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:       createUserDTO.Username,
		Password:       hashedPassword,
		Email:          createUserDTO.Email,
		FullName:       createUserDTO.FullName,
		Role:           createUserDTO.Role,
		ClearanceLevel: createUserDTO.ClearanceLevel,
		IsActive:       true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(id uint, updateUserDTO *dto.UpdateUserDTO) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if updateUserDTO.Password != "" {
		hashedPassword, err := auth.HashPassword(updateUserDTO.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	user.Email = updateUserDTO.Email
	user.FullName = updateUserDTO.FullName
	user.Role = updateUserDTO.Role
	user.ClearanceLevel = updateUserDTO.ClearanceLevel
	user.IsActive = updateUserDTO.IsActive

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
