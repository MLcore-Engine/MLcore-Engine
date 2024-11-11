package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Notebook struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ProjectID       uint      `json:"project_id" gorm:"index;constraint:OnDelete:RESTRICT"`
	Project         Project   `json:"project" gorm:"foreignKey:ProjectID;references:ID"`
	UserID          uint      `json:"user_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	User            User      `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Name            string    `json:"name" gorm:"size:200;unique"`
	Describe        string    `json:"describe" gorm:"size:200"`
	Namespace       string    `json:"namespace" gorm:"size:200;default:jupyter"`
	Image           string    `json:"image" gorm:"size:200;default:''"`
	IDEType         string    `json:"ide_type" gorm:"size:100;default:jupyter"`
	WorkingDir      string    `json:"working_dir" gorm:"size:200;default:''"`
	Env             string    `json:"env" gorm:"size:400;default:''"`
	VolumeMount     string    `json:"volume_mount" gorm:"size:2000"`
	NodeSelector    string    `json:"node_selector" gorm:"size:200;default:notebook=true"`
	ImagePullPolicy string    `json:"image_pull_policy" gorm:"size:20;default:'Always'"`
	ResourceMemory  string    `json:"resource_memory" gorm:"size:100;default:8G"`
	ResourceCPU     string    `json:"resource_cpu" gorm:"size:100;default:4"`
	ResourceGPU     int64     `json:"resource_gpu" gorm:"default:0"`
	Status          string    `json:"status" gorm:"size:50;default:'Creating'"`
	Expand          string    `json:"expand" gorm:"type:text;default:'{}'"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	AccessURL       string    `json:"access_url" gorm:"size:500"`
}

// Insert creates a new Notebook
func (n *Notebook) Insert() error {
	return DB.Create(n).Error
}

// Update updates the Notebook
func (n *Notebook) Update() error {
	return DB.Model(n).Updates(n).Error
}

// Delete removes the Notebook (hard delete)
func (n *Notebook) Delete() error {
	if n.ID == 0 {
		return errors.New("id is empty")
	}
	return DB.Delete(n).Error
}

// GetNotebookByID retrieves a Notebook by ID
func GetNotebookByID(id uint) (*Notebook, error) {
	var notebook Notebook
	result := DB.First(&notebook, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("notebook not found")
		}
		return nil, result.Error
	}
	return &notebook, nil
}

// GetAllNotebooksPaginated retrieves all Notebooks with pagination
func GetAllNotebooksPaginated(offset, limit int) ([]Notebook, int64, error) {
	var notebooks []Notebook
	var total int64

	err := DB.Model(&Notebook{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Offset(offset).Limit(limit).Find(&notebooks).Error
	return notebooks, total, err
}

// GetUserNotebooksPaginated retrieves Notebooks for a specific user with pagination
func GetUserNotebooksPaginated(userID int, offset, limit int) ([]Notebook, int64, error) {
	var notebooks []Notebook
	var total int64

	err := DB.Model(&Notebook{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = DB.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&notebooks).Error
	return notebooks, total, err
}

// SearchNotebooks searches Notebooks by keyword
func SearchNotebooks(keyword string) ([]Notebook, error) {
	var notebooks []Notebook
	err := DB.Where("name LIKE ? OR describe LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&notebooks).Error
	return notebooks, err
}

// Reset updates the Notebook's UpdatedAt timestamp
func (n *Notebook) Reset() error {
	n.UpdatedAt = time.Now()
	return DB.Save(n).Error
}
