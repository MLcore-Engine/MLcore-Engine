package model

import (
	"errors"
	"time"
)

const (
	RoleRoot   = 100
	RoleAdmin  = 10
	RoleCommon = 1
)

type UserProject struct {
	UserID    uint      `json:"userId" gorm:"primaryKey"`
	ProjectID uint      `json:"projectId" gorm:"primaryKey"`
	Role      int       `json:"role" gorm:"type:int;default:1;not null"` // 1000: root, 100: admin, 1: common
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	Project   Project   `json:"-" gorm:"foreignKey:ProjectID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (up *UserProject) ValidateRole() error {
	validRoles := map[int]bool{
		RoleRoot:   true,
		RoleAdmin:  true,
		RoleCommon: true,
	}
	if !validRoles[up.Role] {
		return errors.New("invalid role")
	}
	return nil
}

func (up *UserProject) GetRoleName() string {
	switch up.Role {
	case RoleRoot:
		return "root"
	case RoleAdmin:
		return "admin"
	case RoleCommon:
		return "common"
	default:
		return "unknown"
	}
}