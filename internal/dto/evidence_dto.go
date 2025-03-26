package dto

type CreateTextEvidenceDTO struct {
	CaseID  uint   `json:"caseId" binding:"required"`
	Content string `json:"content" binding:"required"`
	Remarks string `json:"remarks"`
}

type UpdateEvidenceDTO struct {
	Remarks string `json:"remarks"`
}

type EvidenceResponseDTO struct {
	ID        uint   `json:"id"`
	CaseID    uint   `json:"caseId"`
	Type      string `json:"type"`
	Content   string `json:"content,omitempty"`
	ImagePath string `json:"imagePath,omitempty"`
	Remarks   string `json:"remarks,omitempty"`
	AddedBy   string `json:"addedBy"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type DeleteConfirmationDTO struct {
	Confirmation string `json:"confirmation" binding:"required"`
}
