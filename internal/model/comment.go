package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	CaseID    uint   `gorm:"not null"`
	Case      Case   `gorm:"foreignKey:CaseID"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	IPAddress string
}
