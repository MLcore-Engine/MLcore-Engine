package controller

import (
	"MLcore-Engine/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project with the input payload
// @Tags projects
// @Accept json
// @Produce json
// @Param project body model.Project true "Create project"
// @Success 201 {object} model.Project
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /project [post]
func CreateProject(c *gin.Context) {
	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create project: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": project})
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update a project with the input payload
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body model.Project true "Update project"
// @Success 200 {object} model.Project
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /project/{id} [put]
func UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project model.Project
	if err := model.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": project})
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /project/{id} [delete]
func DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := model.DB.Delete(&model.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Project successfully deleted"})
}

// GetProject godoc
// @Summary Get a project
// @Description Get details of a specific project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} model.Project
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /project/{id} [get]
func GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project model.Project
	if err := model.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": project})
}

// ListProjects godoc
// @Summary List projects with pagination
// @Description Get a list of projects with pagination
// @Tags projects
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /project/get-all [get]
func ListProjects(c *gin.Context) {

	fmt.Println("List Projects")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	var projects []model.Project
	var total int64

	if err := model.DB.Model(&model.Project{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to count projects"})
		return
	}

	// 使用 Preload 预加载用户信息
	query := model.DB.Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "display_name", "email") // 选择需要的用户字段
	})

	if err := query.Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to retrieve projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"projects": projects,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}
