package dto

type EvidenceDTO struct {
	CaseID    uint   `json:"caseId" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Content   string `json:"content"`
	ImagePath string `json:"imagePath"`
	Remarks   string `json:"remarks"`
}
