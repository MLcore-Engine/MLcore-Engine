package controller

import (
	"MLcore-Engine/common"
	"MLcore-Engine/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DB预加载辅助函数
func PreloadUserWithoutSensitiveInfo(db *gorm.DB) *gorm.DB {
	return db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "display_name", "role", "status", "email")
	})
}

// GetUserProjects godoc
// @Summary Get projects for a user
// @Description Get all projects associated with a specific user
// @Tags project-members
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} ProjectMembershipResponse
// @Failure 400,500 {object} ProjectMembershipResponse
// @Router /project-members/user/{userId} [get]
func GetUserProjects(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProjectMembershipResponse{
			Success: false,
			Message: "无效的用户ID",
		})
		return
	}

	var user model.User
	if err := model.DB.Preload("Projects").First(&user, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ProjectMembershipResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ProjectMembershipResponse{
		Success: true,
		Data:    user.Projects,
	})
}

// GetProjectMembers godoc
// @Summary Get members of a project
// @Description Get all members of a specific project
// @Tags project-members
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Success 200 {object} []model.User
// @Failure 400,500 {object} ProjectMembershipResponse
// @Router /project-members/project/{projectId} [get]
func GetProjectMembers(c *gin.Context) {
	projectId, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的项目ID",
		})
		return
	}

	var project model.Project
	if err := model.DB.Preload("Users").First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	var formattedUsers []gin.H
	for _, user := range project.Users {
		formattedUsers = append(formattedUsers, gin.H{
			"id":       user.ID,
			"userId":   strconv.FormatUint(uint64(user.ID), 10),
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    formattedUsers,
	})
}

// AddUserToProject godoc
// @Summary Add a user to a project
// @Description Add a user to a project with a specified role
// @Tags project-members
// @Accept json
// @Produce json
// @Param request body SuccessResponse true "Member Info"
// @Success 200 {object} SuccessResponse
// @Failure 400,500 {object} ErrorResponse
// @Router /project-members [post]
func AddUserToProject(c *gin.Context) {
	// 仅使用小驼峰命名
	var request struct {
		UserId    uint `json:"userId" binding:"required"`
		ProjectId uint `json:"projectId" binding:"required"`
		Role      int  `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 记录请求内容
	requestStr := fmt.Sprintf("AddUserToProject 请求: body=%+v", request)
	common.SysLog(requestStr)

	// 验证必要字段
	if request.UserId == 0 || request.ProjectId == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "缺少必要字段: userId和projectId",
		})
		return
	}

	userProject := model.UserProject{
		UserID:    request.UserId,
		ProjectID: request.ProjectId,
		Role:      request.Role,
	}

	if err := userProject.ValidateRole(); err != nil {
		errMsg := fmt.Sprintf("无效的角色值: %d", request.Role)
		common.SysLog(errMsg)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := model.DB.Create(&userProject).Error; err != nil {
		errMsg := fmt.Sprintf("创建用户项目关系失败: %s, userId=%d, projectId=%d", err.Error(), request.UserId, request.ProjectId)
		common.SysLog(errMsg)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// 查询用户信息以返回完整的成员信息
	var user model.User
	if err := model.DB.Select("id", "username").Where("id = ?", request.UserId).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "User found but could not fetch details: " + err.Error(),
		})
		return
	}

	// 返回与前端MemberDTO匹配的结构
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "User successfully added to project",
		Data: MemberDTO{
			ProjectId: userProject.ProjectID,
			UserId:    userProject.UserID,
			Username:  user.Username,
			Role:      userProject.Role,
		},
	})
}

// RemoveUserFromProject godoc
// @Summary Remove a user from a project
// @Description Remove a user from a specified project
// @Tags project-members
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Param userId path int true "User ID"
// @Success 200 {object} ProjectMembershipResponse
// @Failure 400,404,500 {object} ProjectMembershipResponse
// @Router /project-members/{projectId}/{userId} [delete]
func RemoveUserFromProject(c *gin.Context) {
	// 调试日志
	fmt.Printf("请求参数: %+v\n", c.Params)

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProjectMembershipResponse{
			Success: false,
			Message: "无效的用户ID: " + err.Error(),
		})
		return
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProjectMembershipResponse{
			Success: false,
			Message: "无效的项目ID: " + err.Error(),
		})
		return
	}

	result := model.DB.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&model.UserProject{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ProjectMembershipResponse{
			Success: false,
			Message: result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, ProjectMembershipResponse{
			Success: false,
			Message: "用户不在此项目中",
		})
		return
	}

	c.JSON(http.StatusOK, ProjectMembershipResponse{
		Success: true,
		Message: "用户已成功从项目中移除",
		Data:    nil,
	})
}

// UpdateUserProjectRole godoc
// @Summary Update a user's role in a project
// @Description Update the role of a user in a specified project
// @Tags project-members
// @Accept json
// @Produce json
// @Param request body RoleUpdateRequest true "Role Update Info"
// @Success 200 {object} ProjectMembershipResponse
// @Failure 400,404,500 {object} ProjectMembershipResponse
// @Router /project-members [put]
func UpdateUserProjectRole(c *gin.Context) {
	// 仅使用小驼峰命名
	var request struct {
		UserId    uint `json:"userId" binding:"required"`
		ProjectId uint `json:"projectId" binding:"required"`
		Role      int  `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 验证必要字段
	if request.UserId == 0 || request.ProjectId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "缺少必要字段: userId和projectId",
		})
		return
	}

	userProject := model.UserProject{
		UserID:    request.UserId,
		ProjectID: request.ProjectId,
		Role:      request.Role,
	}

	if err := userProject.ValidateRole(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := model.DB.Model(&model.UserProject{}).Where("user_id = ? AND project_id = ?", request.UserId, request.ProjectId).Update("role", request.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 查询用户信息以返回完整的成员信息
	var user model.User
	if err := model.DB.Select("id", "username").Where("id = ?", request.UserId).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "找不到用户详情: " + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户在项目中的角色已成功更新",
		"data": MemberDTO{
			ProjectId: userProject.ProjectID,
			UserId:    userProject.UserID,
			Username:  user.Username,
			Role:      userProject.Role,
		},
	})
}
