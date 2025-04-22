package model

import (
	"time"
)

// model/dataset_entry.go
type DatasetEntry struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	DatasetID   uint      `json:"dataset_id" gorm:"not null;index;constraint:OnDelete:CASCADE"`
	Dataset     Dataset   `json:"-" gorm:"foreignKey:DatasetID;references:ID"`
	EntryIndex  int       `json:"entry_index" gorm:"not null;index"`
	Instruction string    `json:"instruction" gorm:"type:text"`
	Input       string    `json:"input" gorm:"type:text"`
	Output      string    `json:"output" gorm:"type:text"`
	RawContent  string    `json:"-" gorm:"type:text"` // 原始JSON内容
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
