package model

import (
	"gorm.io/gorm"
)

// Base struct for suspect, victim, and witness
type Person struct {
	gorm.Model
	CaseID      uint   `gorm:"not null"`
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Age         int
	Gender      string
	Address     string
	PhoneNumber string
	Notes       string `gorm:"type:text"`
	AddedByID   uint   `gorm:"not null"`
	AddedBy     User   `gorm:"foreignKey:AddedByID"`
}

type Suspect struct {
	Person
	Description string `gorm:"type:text"`
	IsArrested  bool   `gorm:"default:false"`
}

type Victim struct {
	Person
	InjuryDescription string `gorm:"type:text"`
}

type Witness struct {
	Person
	Statement string `gorm:"type:text"`
}
