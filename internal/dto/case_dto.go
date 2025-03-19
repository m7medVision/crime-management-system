package dto

import (
	"github.com/m7medVision/crime-management-system/internal/model"
)

type CaseDTO struct {
	Name               string               `json:"name" binding:"required"`
	Description        string               `json:"description" binding:"required"`
	Area               string               `json:"area" binding:"required"`
	CaseType           string               `json:"caseType" binding:"required"`
	AuthorizationLevel model.ClearanceLevel `json:"authorizationLevel" binding:"required"`
	ReportedByID       uint                 `json:"reportedById" binding:"required"`
}
