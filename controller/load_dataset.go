package controller

import (
	"MLcore-Engine/common"
	"MLcore-Engine/model"
	"MLcore-Engine/services"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImportDataset 导入数据集
// @Summary 导入数据集
// @Description 从JSONL文件导入数据到指定数据集
// @Tags Dataset
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "数据集ID"
// @Param file formData file true "JSONL文件"
// @Success 200 {object} SuccessResponse
// @Router /api/dataset/{id}/import [post]
func ImportDataset(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户ID
	userID, exists := c.Get("id")
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
	if dataset.UserID == userID.(uint) {
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

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "获取上传文件失败: " + err.Error(),
		})
		return
	}

	// 验证文件格式
	ext := filepath.Ext(file.Filename)
	if ext != ".jsonl" && ext != ".json" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "只支持JSONL或JSON文件导入",
		})
		return
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "dataset_import_*.jsonl")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建临时文件失败: " + err.Error(),
		})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 保存上传文件到临时目录
	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存上传文件失败: " + err.Error(),
		})
		return
	}

	// 重新打开文件以读取
	jsonlFile, err := os.Open(tempFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "打开文件失败: " + err.Error(),
		})
		return
	}
	defer jsonlFile.Close()

	// 读取并处理文件内容
	scanner := bufio.NewScanner(jsonlFile)
	lineCount := 0
	validCount := 0

	// 开始数据库事务
	tx := model.DB.Begin()

	// 获取当前最大索引
	var maxIndex struct {
		MaxIndex int
	}
	if dataset.StorageType == "database" || dataset.StorageType == "both" {
		tx.Model(&model.DatasetEntry{}).
			Select("COALESCE(MAX(entry_index), -1) as max_index").
			Where("dataset_id = ?", dataset.ID).
			Scan(&maxIndex)
	}

	currentIndex := maxIndex.MaxIndex + 1

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if line == "" {
			continue
		}

		// 解析JSON行
		var data map[string]string
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			common.SysLog(fmt.Sprintf("解析第%d行失败: %v", lineCount, err))
			continue
		}

		// 检查必要字段
		if data["instruction"] == "" {
			common.SysLog(fmt.Sprintf("第%d行缺少instruction字段", lineCount))
			continue
		}

		// 创建条目
		entry := model.DatasetEntry{
			DatasetID:   dataset.ID,
			EntryIndex:  currentIndex,
			Instruction: data["instruction"],
			Input:       data["input"],
			Output:      data["output"],
			RawContent:  line,
		}

		// 保存到数据库
		if dataset.StorageType == "database" || dataset.StorageType == "both" {
			if err := tx.Create(&entry).Error; err != nil {
				common.SysLog(fmt.Sprintf("保存第%d行到数据库失败: %v", lineCount, err))
				continue
			}
		}

		currentIndex++
		validCount++
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "读取文件出错: " + err.Error(),
		})
		return
	}

	// 如果没有有效行
	if validCount == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "文件中没有有效的数据行",
		})
		return
	}

	// 更新数据集条目计数
	tx.Model(&dataset).Update("entry_count", gorm.Expr("entry_count + ?", validCount))

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存数据失败: " + err.Error(),
		})
		return
	}

	// 如果使用MinIO存储，上传到MinIO
	if dataset.StorageType == "minio" || dataset.StorageType == "both" {
		if dataset.BucketName == "" || dataset.ObjectPath == "" {
			// 配置MinIO存储信息
			bucketName := "datasets"
			objectPath := "dataset_" + id + ".jsonl"

			// 更新数据集
			model.DB.Model(&dataset).Updates(map[string]interface{}{
				"bucket_name": bucketName,
				"object_path": objectPath,
			})

			dataset.BucketName = bucketName
			dataset.ObjectPath = objectPath
		}

		// 将文件上传到MinIO
		jsonlFile.Seek(0, 0) // 重置文件指针到开头
		if err := services.UploadJSONLToMinio(dataset.BucketName, dataset.ObjectPath, jsonlFile); err != nil {
			common.SysLog(fmt.Sprintf("上传到MinIO失败: %v", err))
			// 不中断操作，数据已经保存到数据库
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("成功导入%d条有效数据", validCount),
		"data": map[string]interface{}{
			"total_lines":   lineCount,
			"valid_entries": validCount,
		},
	})
}

// ExportDataset 导出数据集
// @Summary 导出数据集
// @Description 将数据集导出为JSONL文件
// @Tags Dataset
// @Produce application/octet-stream
// @Param id path int true "数据集ID"
// @Success 200
// @Router /api/dataset/{id}/export [get]
func ExportDataset(c *gin.Context) {
	id := c.Param("id")

	// 获取当前用户ID
	userID, exists := c.Get("id")
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
	hasViewPermission := false
	if dataset.UserID == userID.(uint) {
		hasViewPermission = true
	} else if dataset.ProjectID != 0 {
		// 检查是否为项目成员
		var projectUser model.UserProject
		if err := model.DB.Where("project_id = ? AND user_id = ?", dataset.ProjectID, userID).First(&projectUser).Error; err == nil {
			hasViewPermission = true
		}
	}

	if !hasViewPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "您没有权限访问此数据集",
		})
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("dataset_%s_%s.jsonl", id, time.Now().Format("20060102"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")

	// 根据存储类型导出
	if dataset.StorageType == "minio" && dataset.BucketName != "" && dataset.ObjectPath != "" {
		// 从MinIO直接导出
		object, err := services.GetMinioObject(dataset.BucketName, dataset.ObjectPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "从MinIO获取数据失败: " + err.Error(),
			})
			return
		}
		defer object.Close()

		// 直接复制到响应
		if _, err := io.Copy(c.Writer, object); err != nil {
			common.SysLog(fmt.Sprintf("复制MinIO数据到响应失败: %v", err))
			return
		}
	} else {
		// 从数据库导出
		var entries []model.DatasetEntry
		if err := model.DB.Where("dataset_id = ?", id).Order("entry_index ASC").Find(&entries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取数据集条目失败: " + err.Error(),
			})
			return
		}

		// 直接写入响应
		for _, entry := range entries {
			// 如果有原始内容，直接使用
			if entry.RawContent != "" {
				fmt.Fprintln(c.Writer, entry.RawContent)
				continue
			}

			// 否则重新构建JSON
			data := map[string]string{
				"instruction": entry.Instruction,
				"input":       entry.Input,
				"output":      entry.Output,
			}
			jsonBytes, err := json.Marshal(data)
			if err != nil {
				common.SysLog(fmt.Sprintf("序列化条目失败: %v", err))
				continue
			}
			fmt.Fprintln(c.Writer, string(jsonBytes))
		}
	}
}
