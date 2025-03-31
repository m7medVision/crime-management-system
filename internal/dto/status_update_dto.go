package dto

import "github.com/m7medVision/crime-management-system/internal/model"

// StatusUpdateDTO represents the data for updating a case status
type StatusUpdateDTO struct {
	Status model.CaseStatus `json:"status" binding:"required"`
}
