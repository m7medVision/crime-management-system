package model

import (
	"gorm.io/gorm"
)

type EvidenceType string

const (
	EvidenceTypeText  EvidenceType = "text"
	EvidenceTypeImage EvidenceType = "image"
)

type Evidence struct {
	gorm.Model
	CaseID    uint         `gorm:"not null"`
	Case      Case         `gorm:"foreignKey:CaseID"`
	Type      EvidenceType `gorm:"not null"`
	Content   string       `gorm:"type:text"`         // For text evidence
	ImagePath string       `gorm:"type:varchar(255)"` // For image evidence
	Remarks   string       `gorm:"type:text"`
	AddedByID uint         `gorm:"not null"`
	AddedBy   User         `gorm:"foreignKey:AddedByID"`
	IsDeleted bool         `gorm:"default:false"` // For soft delete
}
