package controller

import (
	"MLcore-Engine/model"
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
// @Tags project-memberships
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} ProjectMembershipResponse
// @Failure 400,500 {object} ProjectMembershipResponse
// @Router /project-memberships/user/{userId} [get]
func GetUserProjects(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProjectMembershipResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var user model.User
	if err := model.DB.Preload("Projects").First(&user, userID).Error; err != nil {
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
// @Tags project-memberships
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Success 200 {object} []model.User
// @Failure 400,500 {object} ProjectMembershipResponse
// @Router /project-memberships/project/{projectId} [get]
func GetProjectMembers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid project ID",
		})
		return
	}

	var project model.Project
	if err := model.DB.Preload("Users").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project.Users,
	})
}

// AddUserToProject godoc
// @Summary Add a user to a project
// @Description Add a user to a project with a specified role
// @Tags project-memberships
// @Accept json
// @Produce json
// @Param request body SuccessResponse true "Member Info"
// @Success 200 {object} SuccessResponse
// @Failure 400,500 {object} ErrorResponse
// @Router /project-memberships [post]
func AddUserToProject(c *gin.Context) {

	var request struct {
		UserID    uint `json:"userId"`
		ProjectID uint `json:"projectId"`
		Role      int  `json:"role"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	userProject := model.UserProject{
		UserID:    request.UserID,
		ProjectID: request.ProjectID,
		Role:      request.Role,
	}

	if err := userProject.ValidateRole(); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := model.DB.Create(&userProject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "User successfully added to project",
		Data:    userProject,
	})
}

// RemoveUserFromProject godoc
// @Summary Remove a user from a project
// @Description Remove a user from a specified project
// @Tags project-memberships
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Param userId path int true "User ID"
// @Success 200 {object} ProjectMembershipResponse
// @Failure 400,404,500 {object} ProjectMembershipResponse
// @Router /project-memberships/{projectId}/{userId} [delete]
func RemoveUserFromProject(c *gin.Context) {

	userID, _ := strconv.ParseUint(c.Param("userID"), 10, 32)
	projectID, _ := strconv.ParseUint(c.Param("projectID"), 10, 32)

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
			Message: "User not found in project",
		})
		return
	}

	c.JSON(http.StatusOK, ProjectMembershipResponse{
		Success: true,
		Message: "User successfully removed from project",
		Data:    nil,
	})
}

// UpdateUserProjectRole godoc
// @Summary Update a user's role in a project
// @Description Update the role of a user in a specified project
// @Tags project-memberships
// @Accept json
// @Produce json
// @Param projectId path int true "Project ID"
// @Param userId path int true "User ID"
// @Param request body RoleUpdateRequest true "Role Update Info"
// @Success 200 {object} ProjectMembershipResponse
// @Failure 400,404,500 {object} ProjectMembershipResponse
// @Router /project-memberships [put]
func UpdateUserProjectRole(c *gin.Context) {
	var request struct {
		UserID    uint `json:"userId"`
		ProjectID uint `json:"projectId"`
		Role      int  `json:"role"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	userProject := model.UserProject{
		UserID:    request.UserID,
		ProjectID: request.ProjectID,
		Role:      request.Role,
	}

	if err := userProject.ValidateRole(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := model.DB.Model(&model.UserProject{}).Where("user_id = ? AND project_id = ?", request.UserID, request.ProjectID).Update("role", request.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User role in project successfully updated",
		"data":    userProject,
	})
}

// func validateIDs(c *gin.Context) (projectID uint64, userID uint64, ok bool) {
// 	var err error
// 	projectID, err = strconv.ParseUint(c.Param("projectId"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, ProjectMembershipResponse{
// 			Success: false,
// 			Message: "Invalid project ID",
// 		})
// 		return 0, 0, false
// 	}

// 	userID, err = strconv.ParseUint(c.Param("userId"), 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, ProjectMembershipResponse{
// 			Success: false,
// 			Message: "Invalid user ID",
// 		})
// 		return 0, 0, false
// 	}

// 	return projectID, userID, true
// }
