package dto

type ReportDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location" binding:"required"`
	CivilID     string `json:"civil_id" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Name        string `json:"name" binding:"required"`
}
