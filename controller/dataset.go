package controller

import (
	"MLcore-Engine/common"
	"MLcore-Engine/model"
	"MLcore-Engine/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateDataset 创建新数据集
// @Summary 创建数据集
// @Description 创建新的数据集
// @Tags Dataset
// @Accept json
// @Produce json
// @Param dataset body model.Dataset true "数据集信息"
// @Success 200 {object} DatasetResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/dataset [post]
func CreateDataset(c *gin.Context) {
	var input struct {
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		StorageType      string `json:"storage_type"`
		TemplateType     string `json:"template_type"`
		ProjectID        uint   `json:"project_id"`
		SchemaDefinition string `json:"schema_definition"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 验证项目权限
	if input.ProjectID != 0 {
		var project model.Project
		if err := model.DB.First(&project, input.ProjectID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "项目不存在",
			})
			return
		}

		// 验证用户是否为项目成员
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", input.ProjectID, userID).First(&projectUser).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "您不是该项目的成员，无法关联此项目",
			})
			return
		}
	}

	// 创建数据集
	dataset := model.Dataset{
		Name:             input.Name,
		Description:      input.Description,
		StorageType:      input.StorageType,
		TemplateType:     input.TemplateType,
		ProjectID:        input.ProjectID,
		UserID:           userID.(uint),
		SchemaDefinition: input.SchemaDefinition,
	}

	// 设置默认值
	if dataset.StorageType == "" {
		dataset.StorageType = "database"
	}
	if dataset.TemplateType == "" {
		dataset.TemplateType = "instruction_io"
	}

	// 保存到数据库
	if err := model.DB.Create(&dataset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建数据集失败: " + err.Error(),
		})
		return
	}

	// 如果使用MinIO存储，初始化存储
	if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		bucketName := "datasets"
		objectPath := "dataset_" + dataset.Name + "_" + strconv.FormatUint(uint64(dataset.ID), 10) + ".jsonl"

		// 初始化MinIO存储
		if err := services.InitDatasetMinioStorage(bucketName, objectPath); err != nil {
			// 记录错误但继续处理
			common.SysError(err.Error())
		}

		// 更新数据集存储信息
		dataset.BucketName = bucketName
		dataset.ObjectPath = objectPath
		model.DB.Save(&dataset)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集创建成功",
		"data":    dataset,
	})
}

// ListDatasets 获取数据集列表
// @Summary 获取数据集列表
// @Description 获取当前用户有权限访问的数据集列表
// @Tags Dataset
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} DatasetsResponse
// @Router /api/dataset [get]
func ListDatasets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 获取当前用户ID
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	var datasets []model.Dataset
	var total int64

	// 查询用户有权限访问的数据集
	// 1. 用户创建的数据集
	// 2. 用户所在项目的数据集
	query := model.DB.Model(&model.Dataset{}).
		Joins("LEFT JOIN user_project ON user_project.project_id = datasets.project_id").
		Where("datasets.user_id = ? OR (user_project.user_id = ? AND user_project.deleted_at IS NULL)", userID, userID).
		Group("datasets.id")

	// 计算总数
	query.Count(&total)

	// 获取分页数据
	err := query.
		Preload("User").
		Preload("Project").
		Order("datasets.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&datasets).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取数据集列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": DatasetsListData{
			Datasets: convertToDatasetDTOList(datasets),
			PagedData: PagedData{
				Total: total,
				Page:  page,
				Limit: limit,
			},
		},
	})
}

// GetDataset 获取单个数据集详情
// @Summary 获取数据集详情
// @Description 获取单个数据集的详细信息
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Success 200 {object} DatasetResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/dataset/{id} [get]
func GetDataset(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	var dataset model.Dataset
	err := model.DB.
		Preload("User").
		Preload("Project").
		First(&dataset, id).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证访问权限
	if dataset.UserID != userID.(uint) {
		// 检查是否为项目成员
		if dataset.ProjectID != 0 {
			var projectUser model.UserProject
			if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "您没有权限访问此数据集",
				})
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "您没有权限访问此数据集",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    dataset,
	})
}

// UpdateDataset 更新数据集信息
// @Summary 更新数据集
// @Description 更新数据集的基本信息
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Param dataset body model.Dataset true "数据集信息"
// @Success 200 {object} DatasetResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/dataset/{id} [put]
func UpdateDataset(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 先获取数据集
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证是否为创建者
	if dataset.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "只有数据集创建者可以修改数据集信息",
		})
		return
	}

	// 绑定输入
	var input struct {
		Name             string `json:"name"`
		Description      string `json:"description"`
		ProjectID        uint   `json:"project_id"`
		SchemaDefinition string `json:"schema_definition"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 验证项目权限
	if input.ProjectID != 0 && input.ProjectID != dataset.ProjectID {
		var project model.Project
		if err := model.DB.First(&project, input.ProjectID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "项目不存在",
			})
			return
		}

		// 验证用户是否为项目成员
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", input.ProjectID, userID).First(&projectUser).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "您不是该项目的成员，无法关联此项目",
			})
			return
		}
	}

	// 更新字段
	updates := map[string]interface{}{}

	if input.Name != "" {
		updates["name"] = input.Name
	}

	if input.Description != "" {
		updates["description"] = input.Description
	}

	if input.ProjectID != 0 {
		updates["project_id"] = input.ProjectID
	}

	if input.SchemaDefinition != "" {
		updates["schema_definition"] = input.SchemaDefinition
	}

	// 执行更新
	if err := model.DB.Model(&dataset).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新数据集失败: " + err.Error(),
		})
		return
	}

	// 重新加载数据集信息
	model.DB.
		Preload("User").
		Preload("Project").
		First(&dataset, id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集更新成功",
		"data":    dataset,
	})
}

// DeleteDataset 删除数据集
// @Summary 删除数据集
// @Description 删除指定的数据集
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/dataset/{id} [delete]
func DeleteDataset(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 先获取数据集
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证是否为创建者
	if dataset.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "只有数据集创建者可以删除数据集",
		})
		return
	}

	// 开始事务
	tx := model.DB.Begin()

	// 先删除相关的数据集条目
	if err := tx.Where("dataset_id = ?", id).Delete(&model.DatasetEntry{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除数据集条目失败: " + err.Error(),
		})
		return
	}

	// 删除数据集
	if err := tx.Delete(&dataset).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除数据集失败: " + err.Error(),
		})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除数据集失败: " + err.Error(),
		})
		return
	}

	// 如果使用MinIO存储，尝试删除MinIO中的对象
	if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		if dataset.BucketName != "" && dataset.ObjectPath != "" {
			// 删除MinIO对象
			if err := services.DeleteDatasetMinioObject(dataset.BucketName, dataset.ObjectPath); err != nil {
				// 记录错误但不阻止处理
				common.SysError(err.Error())
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集删除成功",
	})
}

// convertToDatasetDTO 将模型对象转换为DTO
func convertToDatasetDTO(dataset model.Dataset) DatasetDTO {
	return DatasetDTO{
		ID:               dataset.ID,
		Name:             dataset.Name,
		Description:      dataset.Description,
		StorageType:      dataset.StorageType,
		TemplateType:     dataset.TemplateType,
		EntryCount:       dataset.EntryCount,
		TotalSize:        dataset.TotalSize,
		ProjectID:        dataset.ProjectID,
		UserID:           dataset.UserID,
		SchemaDefinition: dataset.SchemaDefinition,
		CreatedAt:        dataset.CreatedAt,
		UpdatedAt:        dataset.UpdatedAt,
	}
}

// convertToDatasetDTOList 将模型对象列表转换为DTO列表
func convertToDatasetDTOList(datasets []model.Dataset) []DatasetDTO {
	dtos := make([]DatasetDTO, len(datasets))
	for i, dataset := range datasets {
		dtos[i] = convertToDatasetDTO(dataset)
	}
	return dtos
}

// SearchDatasets 搜索数据集
// @Summary 搜索数据集
// @Description 根据关键词搜索数据集
// @Tags Dataset
// @Accept json
// @Produce json
// @Param q query string false "搜索关键词"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} DatasetsResponse
// @Router /api/dataset/search [get]
func SearchDatasets(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 获取当前用户ID
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	var datasets []model.Dataset
	var total int64

	// 构造查询
	dbQuery := model.DB.Model(&model.Dataset{}).
		Joins("LEFT JOIN user_project ON user_project.project_id = datasets.project_id").
		Where("datasets.user_id = ? OR (user_project.user_id = ? AND user_project.deleted_at IS NULL)", userID, userID)

	// 添加搜索条件
	if query != "" {
		searchQuery := "%" + query + "%"
		dbQuery = dbQuery.Where("datasets.name LIKE ? OR datasets.description LIKE ?", searchQuery, searchQuery)
	}

	// 去重
	dbQuery = dbQuery.Group("datasets.id")

	// 计算总数
	dbQuery.Count(&total)

	// 获取分页数据
	err := dbQuery.
		Preload("User").
		Preload("Project").
		Order("datasets.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&datasets).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "搜索数据集失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": DatasetsListData{
			Datasets: convertToDatasetDTOList(datasets),
			PagedData: PagedData{
				Total: total,
				Page:  page,
				Limit: limit,
			},
		},
	})
}

// GetProjectDatasets 获取项目的数据集
// @Summary 获取项目数据集
// @Description 获取指定项目的所有数据集
// @Tags Dataset
// @Accept json
// @Produce json
// @Param projectId path int true "项目ID"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Success 200 {object} DatasetsResponse
// @Router /api/dataset/project/{projectId} [get]
func GetProjectDatasets(c *gin.Context) {
	projectID := c.Param("projectId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 获取当前用户ID
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 验证用户是否为项目成员
	var projectUser model.UserProject
	if err := model.DB.Where("project_id = ? AND user_id = ?", projectID, userID).First(&projectUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "您不是该项目的成员，无法访问项目数据集",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "验证项目成员失败: " + err.Error(),
			})
		}
		return
	}

	var datasets []model.Dataset
	var total int64

	// 查询项目数据集
	query := model.DB.Model(&model.Dataset{}).Where("project_id = ?", projectID)

	// 计算总数
	query.Count(&total)

	// 获取分页数据
	err := query.
		Preload("User").
		Preload("Project").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&datasets).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目数据集失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": DatasetsListData{
			Datasets: convertToDatasetDTOList(datasets),
			PagedData: PagedData{
				Total: total,
				Page:  page,
				Limit: limit,
			},
		},
	})
}
