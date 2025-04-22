package model

import (
	"time"
)

type Dataset struct {
	ID               uint   `json:"id" gorm:"primarykey"`
	Name             string `json:"name" gorm:"size:255;not null"`
	Description      string `json:"description" gorm:"type:text"`
	StorageType      string `json:"storage_type" gorm:"size:20;default:'database'"`        // "minio", "database", "both"
	TemplateType     string `json:"template_type" gorm:"size:20;default:'instruction_io'"` // 模板类型
	BucketName       string `json:"bucket_name" gorm:"size:255"`                           // MinIO桶名
	ObjectPath       string `json:"object_path" gorm:"size:255"`                           // MinIO对象路径
	EntryCount       int64  `json:"entry_count" gorm:"default:0"`                          // 条目数量
	TotalSize        int64  `json:"total_size" gorm:"default:0"`                           // 总大小(字节)
	SchemaDefinition string `json:"schema_definition" gorm:"type:text"`                    // JSON Schema定义

	ProjectID uint    `json:"project_id" gorm:"index;constraint:OnDelete:RESTRICT"`
	Project   Project `json:"project" gorm:"foreignKey:ProjectID;references:ID"`

	UserID uint `json:"user_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	User   User `json:"user" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
