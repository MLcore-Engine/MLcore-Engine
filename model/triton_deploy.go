package model

import (
	"errors"

	"gorm.io/gorm"
)

type TritonDeploy struct {
	gorm.Model
	Name      string  `json:"name" gorm:"size:200;unique"`
	Namespace string  `json:"namespace" gorm:"size:200;default:'triton-serving'"`
	Image     string  `json:"image" gorm:"size:200"`
	Replicas  int32   `json:"replicas" gorm:"default:1"`
	CPU       int64   `json:"cpu" gorm:"default:2"`
	Memory    int64   `json:"memory" gorm:"default:4"`
	GPU       int64   `json:"gpu" gorm:"default:0"`
	Status    string  `json:"status" gorm:"size:200;default:'Creating'"`
	Labels    string  `json:"labels" gorm:"type:json"`
	Ports     string  `json:"ports" gorm:"type:json"`
	AccessURL string  `json:"access_url"`
	UserID    uint    `json:"user_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	User      User    `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ProjectID uint    `json:"project_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;references:ID"`

	// Server Configuration
	ModelRepository          string `json:"model_repository" gorm:"default:'/model'"`
	StrictModelConfig        bool   `json:"strict_model_config" gorm:"default:false"`
	AllowPollModelRepository bool   `json:"allow_poll_model_repository" gorm:"default:false"`
	PollRepoSeconds          int    `json:"poll_repo_seconds" gorm:"default:5"`

	// HTTP Configuration
	HttpPort        int  `json:"http_port" gorm:"default:8000"`
	HttpThreadCount int  `json:"http_thread_count" gorm:"default:8"`
	AllowHttp       bool `json:"allow_http" gorm:"default:true"`

	// gRPC Configuration
	GrpcPort                    int  `json:"grpc_port" gorm:"default:8001"`
	GrpcInferAllocationPoolSize int  `json:"grpc_infer_allocation_pool_size" gorm:"default:100"`
	AllowGrpc                   bool `json:"allow_grpc" gorm:"default:true"`

	// Metrics Configuration
	AllowMetrics      bool `json:"allow_metrics" gorm:"default:true"`
	MetricsPort       int  `json:"metrics_port" gorm:"default:8002"`
	MetricsIntervalMs int  `json:"metrics_interval_ms" gorm:"default:2000"`

	// GPU Configuration
	GpuMemoryFraction             float64 `json:"gpu_memory_fraction" gorm:"default:1.0"`
	MinSupportedComputeCapability float64 `json:"min_supported_compute_capability" gorm:"default:6.0"`

	// Logging Configuration
	LogVerbose int  `json:"log_verbose" gorm:"default:0"`
	LogInfo    bool `json:"log_info" gorm:"default:true"`
	LogWarning bool `json:"log_warning" gorm:"default:true"`
	LogError   bool `json:"log_error" gorm:"default:true"`
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
