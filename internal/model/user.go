package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin        Role = "admin"
	RoleInvestigator Role = "investigator"
	RoleOfficer      Role = "officer"
	RoleCitizen      Role = "citizen"
)

type ClearanceLevel string

const (
	ClearanceLow      ClearanceLevel = "low"
	ClearanceMedium   ClearanceLevel = "medium"
	ClearanceHigh     ClearanceLevel = "high"
	ClearanceCritical ClearanceLevel = "critical"
)

type User struct {
	gorm.Model
	Username       string         `gorm:"uniqueIndex;not null"`
	Password       string         `gorm:"not null"`
	Email          string         `gorm:"uniqueIndex;not null"`
	FullName       string         `gorm:"not null"`
	Role           Role           `gorm:"not null;default:'citizen'"`
	ClearanceLevel ClearanceLevel `gorm:"default:'low'"`
	IsActive       bool           `gorm:"default:true"`
	LastLogin      *time.Time
}
