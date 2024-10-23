package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Notebook struct {
	ID        uint `gorm:"primaryKey;comment:id主键"`
	ProjectID uint `gorm:"not null;index;comment:项目组id"` // 定义外键，并创建索引
	// Project          Project        `gorm:"foreignKey:ProjectID;references:ID"`
	UserID          uint           `gorm:"not null;index;comment:用户id"`
	Name            string         `gorm:"size:200;unique;comment:英文名"`
	Describe        string         `gorm:"size:200;comment:描述"`
	Namespace       string         `gorm:"size:200;default:jupyter;comment:命名空间"`
	Images          string         `gorm:"size:200;default:'';comment:镜像"`
	IDEType         string         `gorm:"size:100;default:jupyter;comment:ide类型"`
	WorkingDir      string         `gorm:"size:200;default:'';comment:工作目录"`
	Env             string         `gorm:"size:400;default:'';comment:环境变量"`
	VolumeMount     string         `gorm:"size:2000;default:kubeflow-user-workspace:/mnt,kubeflow-archives:/archives;comment:挂载"`
	NodeSelector    string         `gorm:"size:200;default:cpu=true,notebook=true;comment:机器选择器"`
	ImagePullPolicy string         `gorm:"type:text;size:20;default:'Always';comment:镜像拉取策略"`
	ResourceMemory  string         `gorm:"size:100;default:10G;comment:申请内存"`
	ResourceCPU     string         `gorm:"size:100;default:10;comment:申请cpu"`
	ResourceGPU     int64          `gorm:"size:100;default:0;comment:申请gpu"`
	Status          string         `gorm:"size:50;default:'Creating';comment:notebookStatus"`
	Expand          string         `gorm:"type:text;default:'{}';comment:扩展参数"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	AccessURL       string         `gorm:"size:500;comment:访问URL"`
}

// Insert 创建新的 Notebook
func (n *Notebook) Insert() error {
	return DB.Create(n).Error
}

// Update 更新 Notebook
func (n *Notebook) Update() error {
	// DB.Save()
	return DB.Model(n).Updates(n).Error
}

// Delete 删除 Notebook
func (n *Notebook) Delete() error {
	return DB.Unscoped().Delete(n).Error
}

// GetNotebookByID 根据 ID 获取 Notebook
func GetNotebookByID(id uint) (*Notebook, error) {
	var notebook Notebook
	result := DB.Unscoped().First(&notebook, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("notebook not found")
		}
		return nil, result.Error
	}
	return &notebook, nil
}

// GetAllNotebooksPaginated 获取所有 Notebooks（带分页）
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

// GetUserNotebooksPaginated 获取特定用户的 Notebooks（带分页）
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

// SearchNotebooks 搜索 Notebooks
func SearchNotebooks(keyword string) ([]Notebook, error) {
	var notebooks []Notebook
	err := DB.Where("name LIKE ? OR describe LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&notebooks).Error
	return notebooks, err
}

// Reset 重置 Notebook（这里只是一个示例，具体实现可能需要根据你的业务逻辑来定）
func (n *Notebook) Reset() error {
	n.UpdatedAt = time.Now()
	return DB.Save(n).Error
}
