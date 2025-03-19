package model

import (
	"gorm.io/gorm"
)

type CaseStatus string

const (
	StatusPending CaseStatus = "pending"
	StatusOngoing CaseStatus = "ongoing"
	StatusClosed  CaseStatus = "closed"
)

type Case struct {
	gorm.Model
	CaseNumber         string         `gorm:"uniqueIndex;not null"`
	Name               string         `gorm:"not null"`
	Description        string         `gorm:"type:text;not null"`
	Area               string         `gorm:"not null"` // City/Area
	CaseType           string         `gorm:"not null"`
	Status             CaseStatus     `gorm:"not null;default:'pending'"`
	AuthorizationLevel ClearanceLevel `gorm:"not null;default:'low'"`
	CreatedByID        uint           `gorm:"not null"`
	CreatedBy          User           `gorm:"foreignKey:CreatedByID"`
	ReportedByID       uint           `gorm:"not null"`
	ReportedBy         User           `gorm:"foreignKey:ReportedByID"`
	Reports            []Report       `gorm:"many2many:case_reports;"`
	Assignees          []User         `gorm:"many2many:case_assignees;"`
	Evidence           []Evidence     `gorm:"foreignKey:CaseID"`
	Suspects           []Suspect      `gorm:"foreignKey:CaseID"`
	Victims            []Victim       `gorm:"foreignKey:CaseID"`
	Witnesses          []Witness      `gorm:"foreignKey:CaseID"`
}
