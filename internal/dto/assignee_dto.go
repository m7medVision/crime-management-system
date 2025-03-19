package dto

type AssigneeDTO struct {
	UserID uint `json:"userId" binding:"required"`
}
