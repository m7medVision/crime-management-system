package dto

import "github.com/m7medVision/crime-management-system/internal/model"

type CreateUserDTO struct {
	Username       string               `json:"username" binding:"required"`
	Password       string               `json:"password" binding:"required"`
	Email          string               `json:"email" binding:"required"`
	FullName       string               `json:"fullName" binding:"required"`
	Role           model.Role           `json:"role" binding:"required"`
	ClearanceLevel model.ClearanceLevel `json:"clearanceLevel" binding:"required"`
}

type UpdateUserDTO struct {
	Password       string               `json:"password"`
	Email          string               `json:"email" binding:"required"`
	FullName       string               `json:"fullName" binding:"required"`
	Role           model.Role           `json:"role" binding:"required"`
	ClearanceLevel model.ClearanceLevel `json:"clearanceLevel" binding:"required"`
	IsActive       bool                 `json:"isActive" binding:"required"`
}
