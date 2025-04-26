package controller

import (
	"MLcore-Engine/model"
	"MLcore-Engine/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetDatasetEntries 获取数据集条目列表
// @Summary 获取数据集条目列表
// @Description 获取指定数据集的条目列表
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Param page query int false "页码"
// @Param limit query int false "每页数量"
// @Param q query string false "搜索关键词"
// @Param field query string false "搜索字段(instruction/input/output/all)"
// @Success 200 {object} DatasetEntriesResponse
// @Router /api/dataset/{id}/entries [get]
func GetDatasetEntries(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	query := c.Query("q")
	field := c.DefaultQuery("field", "all")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 获取数据集并验证权限
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证访问权限
	if dataset.UserID != uint(userID.(int)) && dataset.ProjectID != 0 {
		// 检查是否为项目成员
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "您没有权限访问此数据集",
			})
			return
		}
	}

	// 根据存储类型获取数据
	var entries []model.DatasetEntry
	var total int64

	if dataset.StorageType == "database" {
		// 从数据库获取
		dbQuery := model.DB.Model(&model.DatasetEntry{}).Where("dataset_id = ?", id)

		// 添加搜索条件
		if query != "" {
			searchQuery := "%" + query + "%"
			switch field {
			case "instruction":
				dbQuery = dbQuery.Where("instruction LIKE ?", searchQuery)
			case "input":
				dbQuery = dbQuery.Where("input LIKE ?", searchQuery)
			case "output":
				dbQuery = dbQuery.Where("output LIKE ?", searchQuery)
			default:
				dbQuery = dbQuery.Where("instruction LIKE ? OR input LIKE ? OR output LIKE ?",
					searchQuery, searchQuery, searchQuery)
			}
		}

		// 计算总数
		dbQuery.Count(&total)

		// 获取分页数据
		if err := dbQuery.Order("entry_index ASC").Offset(offset).Limit(limit).Find(&entries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取数据集条目失败: " + err.Error(),
			})
			return
		}
	} else if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		// 从MinIO获取
		if dataset.BucketName == "" || dataset.ObjectPath == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "数据集MinIO存储信息不完整",
			})
			return
		}

		// 读取JSONL行
		lines, err := services.ReadJSONLFromMinio(dataset.BucketName, dataset.ObjectPath, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "从MinIO读取数据集条目失败: " + err.Error(),
			})
			return
		}

		// 解析为条目
		entries = make([]model.DatasetEntry, 0, len(lines))
		for i, line := range lines {
			if line == "" {
				continue
			}

			var data map[string]string
			if err := json.Unmarshal([]byte(line), &data); err != nil {
				continue
			}

			entry := model.DatasetEntry{
				DatasetID:   dataset.ID,
				EntryIndex:  offset + i,
				Instruction: data["instruction"],
				Input:       data["input"],
				Output:      data["output"],
				RawContent:  line,
			}

			// 应用搜索过滤
			if query != "" {
				switch field {
				case "instruction":
					if !contains(entry.Instruction, query) {
						continue
					}
				case "input":
					if !contains(entry.Input, query) {
						continue
					}
				case "output":
					if !contains(entry.Output, query) {
						continue
					}
				default:
					if !contains(entry.Instruction, query) &&
						!contains(entry.Input, query) &&
						!contains(entry.Output, query) {
						continue
					}
				}
			}

			entries = append(entries, entry)
		}

		// 估算总数
		// 这里只是一个近似值，如果需要准确值，需要读取全部文件
		if dataset.EntryCount > 0 {
			total = dataset.EntryCount
		} else {
			total = int64(len(entries) + offset)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": DatasetEntriesListData{
			Entries: convertToDatasetEntryDTOList(entries),
			PagedData: PagedData{
				Total: total,
				Page:  page,
				Limit: limit,
			},
		},
	})
}

// GetDatasetEntry 获取单个数据集条目
// @Summary 获取数据集条目
// @Description 获取指定数据集的单个条目
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Param entryId path int true "条目ID或索引"
// @Success 200 {object} DatasetEntryResponse
// @Router /api/dataset/{id}/entry/{entryId} [get]
func GetDatasetEntry(c *gin.Context) {
	id := c.Param("id")
	entryID := c.Param("entryId")

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 获取数据集并验证权限
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证访问权限
	if dataset.UserID != uint(userID.(int)) && dataset.ProjectID != 0 {
		// 检查是否为项目成员
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "您没有权限访问此数据集",
			})
			return
		}
	}

	// 根据存储类型获取条目
	var entry model.DatasetEntry

	if dataset.StorageType == "database" {
		// 尝试按ID获取
		if entryIDInt, err := strconv.Atoi(entryID); err == nil {
			if err := model.DB.Where("dataset_id = ? AND id = ?", id, entryIDInt).First(&entry).Error; err != nil {
				// 如果没有找到，尝试按索引获取
				if err := model.DB.Where("dataset_id = ? AND entry_index = ?", id, entryIDInt).First(&entry).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"success": false,
						"message": "数据集条目不存在",
					})
					return
				}
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的条目ID",
			})
			return
		}
	} else if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		// 从MinIO获取
		if dataset.BucketName == "" || dataset.ObjectPath == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "数据集MinIO存储信息不完整",
			})
			return
		}

		// 转换为索引
		entryIndex, err := strconv.Atoi(entryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "无效的条目索引",
			})
			return
		}

		// 读取单行
		lines, err := services.ReadJSONLFromMinio(dataset.BucketName, dataset.ObjectPath, entryIndex, 1)
		if err != nil || len(lines) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "数据集条目不存在或读取失败",
			})
			return
		}

		line := lines[0]
		var data map[string]string
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "解析数据集条目失败: " + err.Error(),
			})
			return
		}

		entry = model.DatasetEntry{
			DatasetID:   dataset.ID,
			EntryIndex:  entryIndex,
			Instruction: data["instruction"],
			Input:       data["input"],
			Output:      data["output"],
			RawContent:  line,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    entry,
	})
}

// CreateDatasetEntry 创建数据集条目
// @Summary 创建数据集条目
// @Description 向数据集添加新的条目
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Param entry body object true "条目内容"
// @Success 200 {object} DatasetEntryResponse
// @Router /api/dataset/{id}/entry [post]
func CreateDatasetEntry(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 获取数据集并验证权限
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证编辑权限
	hasEditPermission := false
	if dataset.UserID == uint(userID.(int)) {
		hasEditPermission = true
	} else if dataset.ProjectID != 0 {
		// 检查是否为项目成员且有编辑权限
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err == nil {
			// 假设角色值大于1拥有编辑权限
			if projectUser.Role > 1 {
				hasEditPermission = true
			}
		}
	}

	if !hasEditPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "您没有权限编辑此数据集",
		})
		return
	}

	// 绑定输入
	var input struct {
		Instruction string `json:"instruction" binding:"required"`
		Input       string `json:"input"`
		Output      string `json:"output"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 准备新条目
	entry := model.DatasetEntry{
		DatasetID:   dataset.ID,
		Instruction: input.Instruction,
		Input:       input.Input,
		Output:      input.Output,
	}

	// 创建JSON字符串
	jsonData := map[string]string{
		"instruction": input.Instruction,
		"input":       input.Input,
		"output":      input.Output,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成JSON数据失败: " + err.Error(),
		})
		return
	}
	entry.RawContent = string(jsonBytes)

	// 根据存储类型保存
	if dataset.StorageType == "database" || dataset.StorageType == "both" {
		// 获取最大索引
		var maxIndex struct {
			MaxIndex int
		}
		model.DB.Model(&model.DatasetEntry{}).
			Select("COALESCE(MAX(entry_index), -1) as max_index").
			Where("dataset_id = ?", dataset.ID).
			Scan(&maxIndex)

		entry.EntryIndex = maxIndex.MaxIndex + 1

		// 保存到数据库
		if err := model.DB.Create(&entry).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "创建数据集条目失败: " + err.Error(),
			})
			return
		}
	}

	if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		// 保存到MinIO
		if dataset.BucketName == "" || dataset.ObjectPath == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "数据集MinIO存储信息不完整",
			})
			return
		}

		// 追加到JSONL文件
		if err := services.AppendJSONLToMinio(dataset.BucketName, dataset.ObjectPath, string(jsonBytes)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "保存到MinIO失败: " + err.Error(),
			})
			return
		}

		// 如果是仅MinIO模式，则设置条目索引
		if dataset.StorageType == "minio" {
			entry.EntryIndex = int(dataset.EntryCount)
		}
	}

	// 更新数据集条目计数
	model.DB.Model(&dataset).Update("entry_count", gorm.Expr("entry_count + 1"))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集条目创建成功",
		"data":    entry,
	})
}

// UpdateDatasetEntry 更新数据集条目
// @Summary 更新数据集条目
// @Description 更新指定数据集的条目
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Param entryId path int true "条目ID或索引"
// @Param entry body object true "条目内容"
// @Success 200 {object} DatasetEntryResponse
// @Router /api/dataset/{id}/entry/{entryId} [put]
func UpdateDatasetEntry(c *gin.Context) {
	id := c.Param("id")
	entryID := c.Param("entryId")

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 获取数据集并验证权限
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证编辑权限
	hasEditPermission := false
	if dataset.UserID == uint(userID.(int)) {
		hasEditPermission = true
	} else if dataset.ProjectID != 0 {
		// 检查是否为项目成员且有编辑权限
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err == nil {
			// 假设角色值大于1拥有编辑权限
			if projectUser.Role > 1 {
				hasEditPermission = true
			}
		}
	}

	if !hasEditPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "您没有权限编辑此数据集",
		})
		return
	}

	// 绑定输入
	var input struct {
		Instruction string `json:"instruction" binding:"required"`
		Input       string `json:"input"`
		Output      string `json:"output"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 转换为JSON
	jsonData := map[string]string{
		"instruction": input.Instruction,
		"input":       input.Input,
		"output":      input.Output,
	}
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成JSON数据失败: " + err.Error(),
		})
		return
	}
	rawContent := string(jsonBytes)

	// 获取条目索引
	entryIndex, err := strconv.Atoi(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的条目ID",
		})
		return
	}

	// 根据存储类型更新
	if dataset.StorageType == "database" || dataset.StorageType == "both" {
		// 尝试查找现有条目
		var entry model.DatasetEntry
		result := model.DB.Where("dataset_id = ? AND (id = ? OR entry_index = ?)", dataset.ID, entryIndex, entryIndex).First(&entry)

		if result.Error != nil {
			// 条目不存在，创建新条目
			entry = model.DatasetEntry{
				DatasetID:   dataset.ID,
				EntryIndex:  entryIndex,
				Instruction: input.Instruction,
				Input:       input.Input,
				Output:      input.Output,
				RawContent:  rawContent,
			}

			if err := model.DB.Create(&entry).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "创建数据集条目失败: " + err.Error(),
				})
				return
			}
		} else {
			// 更新现有条目
			updates := map[string]interface{}{
				"instruction": input.Instruction,
				"input":       input.Input,
				"output":      input.Output,
				"raw_content": rawContent,
			}

			if err := model.DB.Model(&entry).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "更新数据集条目失败: " + err.Error(),
				})
				return
			}
		}
	}

	if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		// 更新MinIO中的内容
		if dataset.BucketName == "" || dataset.ObjectPath == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "数据集MinIO存储信息不完整",
			})
			return
		}

		// 更新指定行
		if err := services.UpdateJSONLInMinio(dataset.BucketName, dataset.ObjectPath, entryIndex, rawContent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "更新MinIO数据失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集条目更新成功",
		"data": model.DatasetEntry{
			DatasetID:   dataset.ID,
			EntryIndex:  entryIndex,
			Instruction: input.Instruction,
			Input:       input.Input,
			Output:      input.Output,
			RawContent:  rawContent,
		},
	})
}

// DeleteDatasetEntry 删除数据集条目
// @Summary 删除数据集条目
// @Description 删除指定数据集的条目
// @Tags Dataset
// @Accept json
// @Produce json
// @Param id path int true "数据集ID"
// @Param entryId path int true "条目ID或索引"
// @Success 200 {object} SuccessResponse
// @Router /api/dataset/{id}/entry/{entryId} [delete]
func DeleteDatasetEntry(c *gin.Context) {
	id := c.Param("id")
	entryID := c.Param("entryId")

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未授权访问",
		})
		return
	}

	// 获取数据集并验证权限
	var dataset model.Dataset
	if err := model.DB.First(&dataset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "数据集不存在",
		})
		return
	}

	// 验证编辑权限
	hasEditPermission := false
	if dataset.UserID == uint(userID.(int)) {
		hasEditPermission = true
	} else if dataset.ProjectID != 0 {
		// 检查是否为项目成员且有编辑权限
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err == nil {
			// 假设角色值大于1拥有编辑权限
			if projectUser.Role > 1 {
				hasEditPermission = true
			}
		}
	}

	if !hasEditPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "您没有权限编辑此数据集",
		})
		return
	}

	// 获取条目索引
	entryIndex, err := strconv.Atoi(entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的条目ID",
		})
		return
	}

	// 根据存储类型删除
	if dataset.StorageType == "database" || dataset.StorageType == "both" {
		// 删除数据库中的条目
		result := model.DB.Where("dataset_id = ? AND (id = ? OR entry_index = ?)", dataset.ID, entryIndex, entryIndex).Delete(&model.DatasetEntry{})

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "删除数据集条目失败: " + result.Error.Error(),
			})
			return
		}

		if result.RowsAffected == 0 {
			// 如果是both模式，可能只在MinIO中，继续处理
			if dataset.StorageType != "both" {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "数据集条目不存在",
				})
				return
			}
		}
	}

	if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		// MinIO中的删除(用空对象覆盖)
		if dataset.BucketName == "" || dataset.ObjectPath == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "数据集MinIO存储信息不完整",
			})
			return
		}

		// 更新为空对象
		emptyJSON := "{}"
		if err := services.UpdateJSONLInMinio(dataset.BucketName, dataset.ObjectPath, entryIndex, emptyJSON); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "更新MinIO数据失败: " + err.Error(),
			})
			return
		}
	}

	// 更新数据集条目计数
	model.DB.Model(&dataset).Update("entry_count", gorm.Expr("GREATEST(entry_count - 1, 0)"))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据集条目删除成功",
	})
}

// 辅助函数: 检查字符串是否包含子串(不区分大小写)
func contains(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// convertToDatasetEntryDTO 将模型对象转换为DTO
func convertToDatasetEntryDTO(entry model.DatasetEntry) DatasetEntryDTO {
	return DatasetEntryDTO{
		ID:          entry.ID,
		DatasetID:   entry.DatasetID,
		EntryIndex:  entry.EntryIndex,
		Instruction: entry.Instruction,
		Input:       entry.Input,
		Output:      entry.Output,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
	}
}

// convertToDatasetEntryDTOList 将模型对象列表转换为DTO列表
func convertToDatasetEntryDTOList(entries []model.DatasetEntry) []DatasetEntryDTO {
	dtos := make([]DatasetEntryDTO, len(entries))
	for i, entry := range entries {
		dtos[i] = convertToDatasetEntryDTO(entry)
	}
	return dtos
}
