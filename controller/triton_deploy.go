package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"MLcore-Engine/common"
	"MLcore-Engine/model"
	"MLcore-Engine/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
)

// CreateTritonDeploy godoc
// @Summary Create a new Triton Deployment
// @Description Create a new Triton Deployment with the provided details
// @Tags triton_deploy
// @Accept json
// @Produce json
// @Param triton_deploy body model.TritonDeploy true "Triton Deployment details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /triton [post]
func CreateTritonDeploy(c *gin.Context) {

	username := c.GetString("username")
	var deploy model.TritonDeploy
	if err := c.ShouldBindJSON(&deploy); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
			Data:    nil,
		})
		return
	}

	deploy.Name = username + "-tri" + common.GenRandStr(5)
	if deploy.Namespace == "" {
		deploy.Namespace = "triton-serving"
	}
	deploy.Status = "Creating"

	deploy.Labels = "{\"app\":\"" + deploy.Name + "\"}"

	if err := deploy.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to insert TritonDeploy: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Create K8s client
	k8sClient, err := services.NewK8s("services/localconfig")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to create K8s client: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 创建 Triton 配置
	tritonConfig := services.TritonConfig{
		ModelRepository:          deploy.ModelRepository,
		StrictModelConfig:        deploy.StrictModelConfig,
		AllowPollModelRepository: deploy.AllowPollModelRepository,
		PollRepoSeconds:          deploy.PollRepoSeconds,

		HttpPort:        deploy.HttpPort,
		HttpThreadCount: deploy.HttpThreadCount,
		AllowHttp:       deploy.AllowHttp,

		GrpcPort:                    deploy.GrpcPort,
		GrpcInferAllocationPoolSize: deploy.GrpcInferAllocationPoolSize,
		AllowGrpc:                   deploy.AllowGrpc,

		AllowMetrics:      deploy.AllowMetrics,
		MetricsPort:       deploy.MetricsPort,
		MetricsIntervalMs: deploy.MetricsIntervalMs,

		GpuMemoryFraction:             deploy.GpuMemoryFraction,
		MinSupportedComputeCapability: deploy.MinSupportedComputeCapability,

		LogVerbose: deploy.LogVerbose,
		LogInfo:    deploy.LogInfo,
		LogWarning: deploy.LogWarning,
		LogError:   deploy.LogError,
	}

	// 生成部署配置
	deploymentConfig, err := services.GetTritonDeployment(
		deploy.Name,
		deploy.Namespace,
		deploy.Image,
		deploy.Replicas,
		deploy.Labels,
		deploy.CPU,
		deploy.Memory,
		deploy.GPU,
		viper.GetString("triton.mountPath"),
		tritonConfig,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to generate Deployment config: " + err.Error(),
			Data:    nil,
		})
		return
	}

	serviceConfig, err := services.GetTritonService(
		deploy.Name,
		deploy.Namespace,
		deploy.Labels,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to generate Service config: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Create Deployment
	createdDeployment, err := k8sClient.CreateTritonDeployment(deploy.Namespace, deploymentConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to create Deployment: " + err.Error(),
		})
		return
	}

	// Create Service
	createdService, err := k8sClient.CreateTritonService(deploy.Namespace, serviceConfig)
	if err != nil {
		// If Service creation fails, delete the created Deployment
		_ = k8sClient.DeleteDeployment(deploy.Namespace, createdDeployment.Name)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to create Service: " + err.Error(),
		})
		return
	}

	nodePorts := getNodePorts(createdService)
	portsJSON, err := json.Marshal(nodePorts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to marshal node ports: " + err.Error(),
		})
		return
	}
	deploy.Ports = string(portsJSON)
	deploy.AccessURL = fmt.Sprintf("http://%s:%d", viper.GetString("triton.externalIP"), nodePorts[0])
	deploy.Status = "Running"

	if err := deploy.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to update TritonDeploy status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Triton Deployment created successfully",
		"data":    deploy,
	})
}

// UpdateTritonDeploy godoc
// @Summary Update a Triton Deployment
// @Description Update an existing Triton Deployment with the provided details
// @Tags triton_deploy
// @Accept json
// @Produce json
// @Param id path int true "TritonDeploy ID"
// @Param triton_deploy body model.TritonDeploy true "Triton Deployment details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /triton_deploy/{id} [put]
func UpdateTritonDeploy(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid ID parameter",
		})
		return
	}

	deploy, err := model.GetTritonDeployByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var updateData model.TritonDeploy
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid request payload: " + err.Error(),
		})
		return
	}

	// 更新字段
	deploy.Image = updateData.Image
	deploy.Replicas = updateData.Replicas
	deploy.Ports = updateData.Ports
	deploy.CPU = updateData.CPU
	deploy.Memory = updateData.Memory
	deploy.GPU = updateData.GPU

	// 更新 Server 配置
	deploy.ModelRepository = updateData.ModelRepository
	deploy.StrictModelConfig = updateData.StrictModelConfig
	deploy.AllowPollModelRepository = updateData.AllowPollModelRepository
	deploy.PollRepoSeconds = updateData.PollRepoSeconds

	// 更新 HTTP 配置
	deploy.HttpPort = updateData.HttpPort
	deploy.HttpThreadCount = updateData.HttpThreadCount
	deploy.AllowHttp = updateData.AllowHttp

	// 更新 gRPC 配置
	deploy.GrpcPort = updateData.GrpcPort
	deploy.GrpcInferAllocationPoolSize = updateData.GrpcInferAllocationPoolSize
	deploy.AllowGrpc = updateData.AllowGrpc

	// 更新 Metrics 配置
	deploy.AllowMetrics = updateData.AllowMetrics
	deploy.MetricsPort = updateData.MetricsPort
	deploy.MetricsIntervalMs = updateData.MetricsIntervalMs

	// 更新 GPU 配置
	deploy.GpuMemoryFraction = updateData.GpuMemoryFraction
	deploy.MinSupportedComputeCapability = updateData.MinSupportedComputeCapability

	// 更新日志配置
	deploy.LogVerbose = updateData.LogVerbose
	deploy.LogInfo = updateData.LogInfo
	deploy.LogWarning = updateData.LogWarning
	deploy.LogError = updateData.LogError

	if err := deploy.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to update TritonDeploy: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Triton Deployment updated successfully",
		"data":    deploy,
	})
}

// DeleteTritonDeploy godoc
// @Summary Delete a Triton Deployment
// @Description Delete an existing Triton Deployment by ID
// @Tags triton_deploy
// @Accept json
// @Produce json
// @Param id path int true "TritonDeploy ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /triton_deploy/{id} [delete]
func DeleteTritonDeploy(c *gin.Context) {
	// Get ID from params
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Success: false,
			Message: "Invalid ID parameter",
		})
		return
	}

	// Get deployment from database
	deploy, err := model.GetTritonDeployByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Create K8s client
	k8sClient, err := services.NewK8s("services/localconfig")
	if err != nil {
		common.SysError(err.Error())
	} else {
		// Delete K8s Deployment
		if err := k8sClient.DeleteDeployment(deploy.Namespace, deploy.Name); err != nil {
			common.SysError(err.Error())
		}

		// Delete K8s Service
		if err := k8sClient.DeleteService2(deploy.Namespace, deploy.Name); err != nil {
			common.SysError(err.Error())
		}
	}

	// Update status and soft delete from database
	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update status to Deleted
	if err := tx.Model(&deploy).Update("status", "Deleted").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to update deployment status: " + err.Error(),
		})
		return
	}

	// Soft delete the record
	if err := tx.Delete(&deploy).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to delete deployment from database: " + err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to commit transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Triton Deployment deleted successfully",
	})
}

// ListTritonDeploys godoc
// @Summary List Triton Deployments
// @Description Get a paginated list of Triton Deployments
// @Tags triton_deploy
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} TritonDeployResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /triton/get-all [get]
func ListTritonDeploys(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Get user information
	// username := c.GetString("username")
	userID := c.GetInt("user_id")
	fmt.Println(userID)
	role := c.GetInt("role")

	var deploys []model.TritonDeploy
	var total int64
	query := model.DB.Model(&model.TritonDeploy{})

	// Filter by username if not root user
	if role != 100 {
		query = query.Where("user_id = ?", userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to count Triton deployments: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Get paginated deployments
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&deploys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Message: "Failed to retrieve Triton deployments: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// Process ports for each deployment
	for i := range deploys {
		var ports []int32
		if deploys[i].Ports != "" {
			if err := json.Unmarshal([]byte(deploys[i].Ports), &ports); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Success: false,
					Message: "Failed to parse ports data: " + err.Error(),
					Data:    nil,
				})
				return
			}
		}
		// You can add the parsed ports to a new field if needed
		// deploys[i].ParsedPorts = ports
	}

	c.JSON(http.StatusOK, TritonDeployResponse{
		Success: true,
		Message: "",
		Data: TritonDeployListData{
			Deployments: deploys,
			Total:       total,
			Page:        page,
			Limit:       limit,
		},
	})
}

// GetTritonConfig godoc
// @Summary Get Triton configuration from config file
// @Description Get all available Triton configuration options from config.yaml
// @Tags triton
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/triton/config [get]
func GetTritonConfig(c *gin.Context) {
	config := map[string]interface{}{
		"namespace":  viper.GetString("triton.namespace"),
		"externalIP": viper.GetString("triton.externalIP"),
		"images":     viper.GetStringSlice("triton.images"),
		"resources": map[string]interface{}{
			"cpu":    viper.GetStringMap("triton.resources.cpu"),
			"memory": viper.GetStringMap("triton.resources.memory"),
			"gpu":    viper.GetStringMap("triton.resources.gpu"),
		},
		"ports": map[string]interface{}{
			"http":    viper.GetIntSlice("triton.ports.http"),
			"grpc":    viper.GetIntSlice("triton.ports.grpc"),
			"metrics": viper.GetIntSlice("triton.ports.metrics"),
		},
		"backend":          viper.GetStringSlice("triton.backend"),
		"model_repository": viper.GetStringSlice("triton.model_repository"),
		"logging": map[string]interface{}{
			"verbose": viper.GetIntSlice("triton.logging.verbose"),
			"formats": viper.GetStringSlice("triton.logging.formats"),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Triton configuration retrieved successfully",
		"data":    config,
	})
}

// Helper function to get NodePort from Service
func getNodePorts(service *corev1.Service) []int32 {
	var ports []int32
	for _, port := range service.Spec.Ports {
		if port.NodePort != 0 {
			ports = append(ports, port.NodePort)
		}
	}
	return ports
}

// Response structs
type TritonDeployResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:""`
	Data    TritonDeployListData `json:"data"`
}

type TritonDeployListData struct {
	Deployments []model.TritonDeploy `json:"deployments"`
	Total       int64                `json:"total" example:"10"`
	Page        int                  `json:"page" example:"1"`
	Limit       int                  `json:"limit" example:"20"`
}
