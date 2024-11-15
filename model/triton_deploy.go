package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TritonDeploy struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	User         User           `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ProjectID    uint           `json:"project_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	Project      Project        `json:"project" gorm:"foreignKey:ProjectID;references:ID"`
	Name         string         `json:"name" gorm:"size:200;unique"`
	Image        string         `json:"image" gorm:"size:200;default:'nvcr.io/nvidia/tritonserver:20.12-py3'"`
	Replicas     int32          `json:"replicas" gorm:"default:1"`
	Ports        string         `json:"ports" gorm:"type:json"`         // JSON-encoded ports
	Labels       string         `json:"labels" gorm:"type:json"`        // JSON-encoded labels
	VolumeMounts string         `json:"volume_mounts" gorm:"type:json"` // JSON-encoded volume mounts
	Resources    string         `json:"resources" gorm:"type:json"`     // JSON-encoded resources
	Command      string         `json:"command" gorm:"type:text"`       // Command to run
	Args         string         `json:"args" gorm:"type:text"`          // Arguments for the command
	CPU          int64          `json:"cpu" gorm:"default:2"`
	Memory       int64          `json:"memory" gorm:"default:4"`
	GPU          int64          `json:"gpu" gorm:"default:0"`
	Namespace    string         `json:"namespace" gorm:"size:200;default:'triton-serving'"`
	Status       string         `json:"status" gorm:"size:200;default:'Creating'"`
	AccessURL    string         `json:"access_url"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Insert creates a new TritonDeploy
func (t *TritonDeploy) Insert() error {
	return DB.Create(t).Error
}

// Update updates the TritonDeploy
func (t *TritonDeploy) Update() error {
	return DB.Model(t).Updates(t).Error
}

// Delete removes the TritonDeploy (soft delete)
func (t *TritonDeploy) Delete() error {
	if t.ID == 0 {
		return errors.New("id is empty")
	}
	return DB.Delete(t).Error
}

// GetTritonDeployByID retrieves a TritonDeploy by ID
func GetTritonDeployByID(id uint) (*TritonDeploy, error) {
	var deploy TritonDeploy
	result := DB.First(&deploy, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("triton deploy not found")
		}
		return nil, result.Error
	}
	return &deploy, nil
}

// GetAllTritonDeploysPaginated retrieves all TritonDeploys with pagination
func GetAllTritonDeploysPaginated(offset, limit int) ([]TritonDeploy, int64, error) {
	var deploys []TritonDeploy
	var total int64

	err := DB.Model(&TritonDeploy{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Offset(offset).Limit(limit).Find(&deploys).Error
	return deploys, total, err
}

// GetUserTritonDeploysPaginated retrieves TritonDeploys for a specific user with pagination
func GetUserTritonDeploysPaginated(userID int, offset, limit int) ([]TritonDeploy, int64, error) {
	var deploys []TritonDeploy
	var total int64

	err := DB.Model(&TritonDeploy{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&deploys).Error
	return deploys, total, err
}

// SearchTritonDeploys searches TritonDeploys by keyword
func SearchTritonDeploys(keyword string) ([]TritonDeploy, error) {
	var deploys []TritonDeploy
	err := DB.Where("name LIKE ?", "%"+keyword+"%").Find(&deploys).Error
	return deploys, err
}
