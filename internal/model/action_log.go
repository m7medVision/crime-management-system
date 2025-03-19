package model

import (
	"gorm.io/gorm"
)

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

type AuditLog struct {
	gorm.Model
	UserID     uint       `gorm:"not null"`
	User       User       `gorm:"foreignKey:UserID"`
	Action     ActionType `gorm:"not null"`
	EntityType string     `gorm:"not null"` // "evidence", "case", "user", etc.
	EntityID   uint       `gorm:"not null"`
	OldValue   string     `gorm:"type:text"`
	NewValue   string     `gorm:"type:text"`
	IPAddress  string
}
