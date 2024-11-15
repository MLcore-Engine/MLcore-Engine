package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type TrainingJob struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	User            User           `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ProjectID       uint           `json:"project_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	Project         Project        `json:"project" gorm:"foreignKey:ProjectID;references:ID"`
	Name            string         `json:"name" gorm:"size:200;unique"`
	Parameters      string         `json:"parameters" gorm:"type:text"` // JSON string for parameters
	Image           string         `json:"image" gorm:"size:200;default:'''"`
	ImagePullPolicy string         `json:"image_pull_policy" gorm:"size:200;default:'IfNotPresent'"`
	Status          string         `json:"status" gorm:"size:50;default:'Pending'"`
	Namespace       string         `json:"namespace" gorm:"size:200;default:'train'"`
	RestartPolicy   string         `json:"restart_policy" gorm:"size:200;default:'OnFailure'"`
	Command         string         `json:"command" gorm:"type:text"`        // JSON-encoded array of commands
	Args            string         `json:"args,omitempty" gorm:"type:text"` // JSON-encoded array of arguments
	MasterReplicas  int32          `json:"master_replicas"`                 // Number of master replicas
	WorkerReplicas  int32          `json:"worker_replicas"`                 // Number of worker replicas
	GPUsPerNode     int64          `json:"gpus_per_node"`                   // Number of GPUs per node
	CPULimit        string         `json:"cpu_limit"`                       // CPU limit per container
	MemoryLimit     string         `json:"memory_limit"`                    // Memory limit per container
	NodeSelector    string         `json:"node_selector" gorm:"type:json"`  // JSON-encoded node selector
	Env             string         `json:"env" gorm:"type:json"`            // JSON-encoded environment variables
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// Insert creates a new TrainingJob
func (t *TrainingJob) Insert() error {
	return DB.Create(t).Error
}

// Update updates the TrainingJob
func (t *TrainingJob) Update() error {
	return DB.Model(t).Updates(t).Error
}

// Delete removes the TrainingJob (soft delete)
func (t *TrainingJob) Delete() error {
	if t.ID == 0 {
		return errors.New("id is empty")
	}
	return DB.Delete(t).Error
}

// GetTrainingJobByID retrieves a TrainingJob by ID
func GetTrainingJobByID(id uint) (*TrainingJob, error) {
	var job TrainingJob
	result := DB.First(&job, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("training job not found")
		}
		return nil, result.Error
	}
	return &job, nil
}

// GetAllTrainingJobsPaginated retrieves all TrainingJobs with pagination
func GetAllTrainingJobsPaginated(offset, limit int) ([]TrainingJob, int64, error) {
	var jobs []TrainingJob
	var total int64

	err := DB.Model(&TrainingJob{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Offset(offset).Limit(limit).Find(&jobs).Error
	return jobs, total, err
}

// GetUserTrainingJobsPaginated retrieves TrainingJobs for a specific user with pagination
func GetUserTrainingJobsPaginated(userID int, offset, limit int) ([]TrainingJob, int64, error) {
	var jobs []TrainingJob
	var total int64

	err := DB.Model(&TrainingJob{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&jobs).Error
	return jobs, total, err
}

// SearchTrainingJobs searches TrainingJobs by keyword
func SearchTrainingJobs(keyword string) ([]TrainingJob, error) {
	var jobs []TrainingJob
	err := DB.Where("name LIKE ? OR framework LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&jobs).Error
	return jobs, err
}
