package dto

type ReportDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Location    string `json:"location" binding:"required"`
	ReportedByID uint  `json:"reportedById" binding:"required"`
}
