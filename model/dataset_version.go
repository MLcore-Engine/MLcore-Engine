package model

import (
	"time"
)

// model/dataset_version.go
type DatasetVersion struct {
	ID          uint    `json:"id" gorm:"primarykey"`
	DatasetID   uint    `json:"dataset_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Dataset     Dataset `json:"-" gorm:"foreignKey:DatasetID;references:ID"`
	Version     string  `json:"version" gorm:"size:50;not null"`
	Description string  `json:"description" gorm:"type:text"`
	IsActive    bool    `json:"is_active" gorm:"default:false"`
	BucketName  string  `json:"bucket_name" gorm:"size:255"`
	ObjectPath  string  `json:"object_path" gorm:"size:255"`
	EntryCount  int64   `json:"entry_count" gorm:"default:0"`

	// 创建者关联
	UserID uint `json:"user_id" gorm:"not null;index;constraint:OnDelete:RESTRICT"`
	User   User `json:"-" gorm:"foreignKey:UserID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
}
