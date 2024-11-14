package controller

import (
	"MLcore-Engine/common"
	"MLcore-Engine/model"
	"MLcore-Engine/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// CreateTrainingJob godoc
// @Summary Create a new Training Job
// @Description Create a new Training Job with the provided details
// @Tags training
// @Accept json
// @Produce json
// @Param training_job body model.TrainingJob true "Training Job details"
// @Success 200 {object} TrainingJobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pytorchtrain [post]
func CreateTrainingJob(c *gin.Context) {

	fmt.Println("CreateTrainingJob")
	username := c.GetString("username")
	var job model.TrainingJob
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload: " + err.Error(),
			"data":    nil,
		})
		return
	}

	job.Name = username + "-pytorchjob-" + common.GenRandStr(5)
	job.Status = "Pending"

	// Insert TrainingJob into the database
	if err := job.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to insert Training Job: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Create Kubernetes client
	k8sClient, err := services.NewK8s("./services/localconfig")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create K8s client: " + err.Error(),
			"data":    nil,
		})
		// Optionally, rollback the database insertion if K8s client creation fails
		if delErr := job.Delete(); delErr != nil {
			fmt.Printf("Failed to rollback Training Job: %v\n", delErr)
		}
		return
	}

	// Create Kubernetes PyTorchJob
	err = createPyTorchJob(k8sClient, &job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create PyTorch Job: " + err.Error(),
			"data":    nil,
		})
		// Optionally, rollback the database insertion if K8s Job creation fails
		if delErr := job.Delete(); delErr != nil {
			fmt.Printf("Failed to rollback Training Job: %v\n", delErr)
		}
		return
	}

	// Update TrainingJob status to 'Running'
	job.Status = "Running"
	if err := job.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update Training Job status: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, TrainingJobResponse{
		Success: true,
		Message: "Training Job created successfully",
		Data:    job,
	})
}

// createPyTorchJob creates a Kubernetes PyTorchJob and updates the TrainingJob with Kubernetes details
func createPyTorchJob(k8sClient *services.K8s, job *model.TrainingJob) error {
	// set default values
	if job.Image == "" {
		return fmt.Errorf("image is empty: %v", job.Image)
	}

	// parse Command string to slice
	var command []string
	// if err := json.Unmarshal([]byte(job.Command), &command); err != nil || len(command) == 0 {
	// 	command = []string{}
	// }

	// parse Args string to slice (optional)
	var args []string
	if job.Args != "" {
		if err := json.Unmarshal([]byte(job.Args), &args); err != nil {
			return fmt.Errorf("invalid args format: %v", err)
		}
	}

	// set default replicas
	if job.MasterReplicas == 0 {
		job.MasterReplicas = 1
	}
	if job.WorkerReplicas == 0 {
		job.WorkerReplicas = 1
	}

	if job.CPULimit == "" {
		job.CPULimit = "4"
	}
	if job.MemoryLimit == "" {
		job.MemoryLimit = "8Gi"
	}

	// parse NodeSelector string to map
	nodeSelector := make(map[string]string)
	// if job.NodeSelector != "" {
	// 	if err := json.Unmarshal([]byte(job.NodeSelector), &nodeSelector); err != nil {
	// 		return fmt.Errorf("invalid node selector format: %v", err)
	// 	}
	// }

	// parse Env string to EnvVar slice
	var envVars []services.EnvVar
	// if job.Env != "" {
	// 	if err := json.Unmarshal([]byte(job.Env), &envVars); err != nil {
	// 		return fmt.Errorf("invalid env vars format: %v", err)
	// 	}
	// }

	if job.RestartPolicy == "" {
		job.RestartPolicy = "OnFailure"
	}

	if job.ImagePullPolicy == "" {
		job.ImagePullPolicy = "IfNotPresent"
	}

	// 创建 PyTorchJob 配置
	config := services.PyTorchJobConfig{
		Name:            job.Name,
		Namespace:       job.Namespace,
		Image:           job.Image,
		ImagePullPolicy: job.ImagePullPolicy,
		RestartPolicy:   job.RestartPolicy,
		Command:         command,
		Args:            args,
		MasterReplicas:  job.MasterReplicas,
		WorkerReplicas:  job.WorkerReplicas,
		GPUsPerNode:     job.GPUsPerNode,
		CPULimit:        job.CPULimit,
		MemoryLimit:     job.MemoryLimit,
		NodeSelector:    nodeSelector,
		Env:             envVars,
	}

	// 在 Kubernetes 中创建 PyTorchJob
	_, err := k8sClient.CreatePyTorchJob(config.Namespace, config)
	if err != nil {
		return fmt.Errorf("failed to create PyTorchJob: %v", err)
	}

	return nil
}

// interfaceSliceToStringSlice converts []interface{} to []string
// func interfaceSliceToStringSlice(slice []interface{}) []string {
// 	strSlice := make([]string, len(slice))
// 	for i, v := range slice {
// 		strSlice[i] = fmt.Sprintf("%v", v)
// 	}
// 	return strSlice
// }

// DeleteTrainingJob godoc
// @Summary Delete a Training Job
// @Description Delete a Training Job by its ID
// @Tags training
// @Accept json
// @Produce json
// @Param id path int true "Training Job ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pytorchtrain/{id} [delete]
func DeleteTrainingJob(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid id parameter",
			"data":    nil,
		})
		return
	}

	job, err := model.GetTrainingJobByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Training Job not found",
			"data":    nil,
		})
		return
	}

	// Create Kubernetes client
	k8sClient, err := services.NewK8s("./services/localconfig")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create K8s client: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// Delete Kubernetes resources based on framework
	err = k8sClient.DeletePyTorchJob(job.Namespace, job.Name)
	if err != nil && !k8serrors.IsNotFound(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete PyTorch Job: " + err.Error(),
		})
		return
	}

	// Delete TrainingJob from the database
	if err := job.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete Training Job from database: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Training Job deleted successfully",
	})
}

// GetTrainingJob godoc
// @Summary Get a Training Job
// @Description Get a Training Job by its ID
// @Tags training
// @Accept json
// @Produce json
// @Param id path int true "Training Job ID"
// @Success 200 {object} TrainingJobResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /training/{id} [get]
func GetTrainingJob(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid id parameter",
		})
		return
	}

	job, err := model.GetTrainingJobByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Training Job not found",
		})
		return
	}

	c.JSON(http.StatusOK, TrainingJobResponse{
		Success: true,
		Message: "",
		Data:    *job,
	})
}

// ListTrainingJobs godoc
// @Summary List Training Jobs
// @Description Get a paginated list of Training Jobs
// @Tags training
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} TrainingJobsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /training/get-all [get]
func ListTrainingJobs(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Get user role and ID
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not authenticated"})
		return
	}

	userId, _ := c.Get("user_id")

	var jobs []model.TrainingJob
	var total int64
	query := model.DB.Model(&model.TrainingJob{}).Preload("User")

	// Filter by user if not admin
	if role.(int) != 100 {
		query = query.Where("user_id = ?", userId)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to count training jobs", "data": nil})
		return
	}

	// Get paginated jobs
	if err := query.Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to retrieve training jobs", "data": nil})
		return
	}

	c.JSON(http.StatusOK, TrainingJobsResponse{
		Success: true,
		Message: "",
		Data: TrainingJobsListData{
			TrainingJobs: jobs,
			Total:        total,
			Page:         page,
			Limit:        limit,
		},
	})
}

// Define response structs
type TrainingJobResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Training Job created successfully"`
	Data    model.TrainingJob `json:"data"`
}

type TrainingJobsResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:""`
	Data    TrainingJobsListData `json:"data"`
}

type TrainingJobsListData struct {
	TrainingJobs []model.TrainingJob `json:"training_jobs"`
	Total        int64               `json:"total" example:"10"`
	Page         int                 `json:"page" example:"1"`
	Limit        int                 `json:"limit" example:"20"`
}

type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation successful"`
}
