package dto

type UserDTO struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	FullName       string `json:"fullName" binding:"required"`
	Password       string `json:"password,omitempty"`
	Role           string `json:"role" binding:"required,oneof=admin investigator officer citizen"`
	ClearanceLevel string `json:"clearanceLevel" binding:"required,oneof=low medium high critical"`
	IsActive       bool   `json:"isActive"`
}

type UpdateUserDTO struct {
	Email          string `json:"email,omitempty" binding:"omitempty,email"`
	FullName       string `json:"fullName,omitempty"`
	Password       string `json:"password,omitempty"`
	Role           string `json:"role,omitempty" binding:"omitempty,oneof=admin investigator officer citizen"`
	ClearanceLevel string `json:"clearanceLevel,omitempty" binding:"omitempty,oneof=low medium high critical"`
	IsActive       *bool  `json:"isActive,omitempty"`
}
