package model

import (
	"gorm.io/gorm"
)

// Project represents a project in the system
type Project struct {
	gorm.Model
	Name        string     `json:"name" gorm:"uniqueIndex;not null" validate:"required,max=50"`
	Description string     `json:"description" validate:"max=500"`
	Users       []User     `json:"users" gorm:"many2many:user_projects"`
	Notebooks   []Notebook `json:"notebooks" gorm:"foreignKey:ProjectID;constraint:OnDelete:RESTRICT"`
}
