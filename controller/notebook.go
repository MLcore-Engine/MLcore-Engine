package controller

import (
	"MLcore-Engine/common"
	"MLcore-Engine/model"
	"MLcore-Engine/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// CreateNotebook godoc
// @Summary Create a new Notebook
// @Description Create a new Notebook with the provided details
// @Tags notebook
// @Accept json
// @Produce json
// @Param notebook body model.Notebook true "Notebook details"
// @Success 200 {object} NotebookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notebook [post]
func CreateNotebook(c *gin.Context) {
	username := c.GetString("username")
	var notebook model.Notebook
	if err := c.ShouldBindJSON(&notebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload: " + err.Error(),
		})
		return
	}

	// Set Notebook basic information
	setNotebookDefaults(&notebook, username)

	// Insert Notebook into database
	if err := notebook.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to insert Notebook: " + err.Error(),
		})
		return
	}

	// Create K8s client
	k8sClient, err := services.NewK8s("./services/config")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	labels := map[string]string{
		"app":      notebook.Name,
		"pod-type": "notebook",
		"user":     username,
	}

	// Create Pod
	createdPod, err := createPodForNotebook(k8sClient, &notebook, username, labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create Pod: " + err.Error(),
		})
		return
	}

	// Create Service
	_, err = createServiceForNotebook(k8sClient, &notebook, labels)
	if err != nil {
		// If Service creation fails, delete the created Pod
		_ = k8sClient.DeletePod(notebook.Namespace, createdPod.Name)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create Service: " + err.Error(),
		})
		return
	}

	_, err = createVsForNotebook(k8sClient, &notebook, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create VirtualService: " + err.Error(),
		})
		return
	}

	// Update Notebook status
	notebook.Status = "Creating"
	notebook.Name = createdPod.Name

	if err := notebook.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update Notebook status: " + err.Error(),
		})
		return
	}

	simplifiedPodInfo := gin.H{
		"Id":             notebook.ID,
		"Name":           createdPod.Name,
		"Describe":       fmt.Sprintf("Notebook for user %s", username),
		"ResourceMemory": notebook.ResourceMemory,
		"ResourceCPU":    notebook.ResourceCPU,
		"ResourceGPU":    notebook.ResourceGPU,
		"Status":         notebook.Status,
		"AccessURL":      notebook.AccessURL,
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notebook and Service created successfully",
		"data": gin.H{
			"pod": simplifiedPodInfo,
		},
	})
}

// setNotebookDefaults sets default values for a Notebook
func setNotebookDefaults(notebook *model.Notebook, username string) {
	notebook.Namespace = viper.GetString("notebook.namespace")
	notebook.Name = username
	notebook.AccessURL = fmt.Sprintf("http://%s/notebook/%s/%s/lab?#%s",
		viper.GetString("notebook.externalIP"), notebook.Namespace, notebook.Name, "mnt/"+username)
}

// createPodForNotebook creates a Pod for the Notebook
func createPodForNotebook(k8sClient *services.K8s, notebook *model.Notebook, username string, labels map[string]string) (*corev1.Pod, error) {
	command := []string{"sh", "-c"}
	argTemplate := `jupyter lab --notebook-dir=/mnt/%s --ip=0.0.0.0 --no-browser --allow-root --port=3000 --NotebookApp.token='' --NotebookApp.password='' --ServerApp.disable_check_xsrf=True --NotebookApp.allow_origin='*' --NotebookApp.base_url=/notebook/jupyter/%s/ --NotebookApp.tornado_settings='{"headers": {"Content-Security-Policy": "frame-ancestors * 'self' "}}'`
	args := []string{fmt.Sprintf(argTemplate, username, username)}

	userWorkspaceVolume := viper.GetString("notebook.volumes.userWorkspace")
	archivesVolume := viper.GetString("notebook.volumes.archives")

	volumeMounts := []corev1.VolumeMount{
		{Name: userWorkspaceVolume, MountPath: fmt.Sprintf("/mnt/%s", username), SubPath: username},
		{Name: archivesVolume, MountPath: fmt.Sprintf("/archives/%s", username), SubPath: username},
		{Name: "tz-config", MountPath: "/etc/localtime"},
		{Name: "dshm", MountPath: "/dev/shm"},
	}

	volumes := []corev1.Volume{
		{Name: userWorkspaceVolume, VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: userWorkspaceVolume}}},
		{Name: archivesVolume, VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: archivesVolume}}},
		{Name: "tz-config", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/etc/localtime"}}},
		{Name: "dshm", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: corev1.StorageMediumMemory}}},
	}

	env := []corev1.EnvVar{
		{Name: "NO_AUTH", Value: "true"},
		{Name: "USERNAME", Value: username},
		{Name: "NODE_OPTIONS", Value: "--max-old-space-size=4096"},
		{Name: "K8S_NODE_NAME", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "spec.nodeName"}}},
		{Name: "K8S_POD_NAMESPACE", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"}}},
		{Name: "K8S_POD_IP", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}}},
		{Name: "K8S_HOST_IP", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"}}},
		{Name: "K8S_POD_NAME", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}}},
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        notebook.Name,
			Namespace:   notebook.Namespace,
			Labels:      labels,
			Annotations: make(map[string]string),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            notebook.Name,
					Image:           viper.GetString("notebook.image.notebookcpu"),
					Command:         command,
					Args:            args,
					WorkingDir:      fmt.Sprintf("/mnt/%s", username),
					VolumeMounts:    volumeMounts,
					Env:             env,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse(notebook.ResourceMemory),
							corev1.ResourceCPU:    resource.MustParse(notebook.ResourceCPU),
							"nvidia.com/gpu":      *resource.NewQuantity(notebook.ResourceGPU, resource.DecimalSI),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse(notebook.ResourceMemory),
							corev1.ResourceCPU:    resource.MustParse(notebook.ResourceCPU),
							"nvidia.com/gpu":      *resource.NewQuantity(notebook.ResourceGPU, resource.DecimalSI),
						},
					},
				},
			},
			Volumes:            volumes,
			RestartPolicy:      corev1.RestartPolicyNever,
			NodeSelector:       map[string]string{"notebook": "true"},
			ImagePullSecrets:   []corev1.LocalObjectReference{{Name: "hubsecret"}},
			ServiceAccountName: "default",
			SchedulerName:      viper.GetString("notebook.schedule"),
		},
	}
	pod.ObjectMeta.ResourceVersion = ""
	pod.ObjectMeta.UID = ""
	pod.ObjectMeta.CreationTimestamp = metav1.Time{}

	return k8sClient.CreatePod(notebook.Namespace, pod)
}

// createServiceForNotebook creates a Service for the Notebook
func createServiceForNotebook(k8sClient *services.K8s, notebook *model.Notebook, labels map[string]string) (*corev1.Service, error) {
	port := viper.GetInt("notebook.defaultPort")
	return k8sClient.CreateServiceForNotebook(notebook.Namespace, notebook.Name, port, labels)
}

// createVsForNotebook creates a VirtualService for the Notebook
func createVsForNotebook(k8sClient *services.K8s, notebook *model.Notebook, username string) (*unstructured.Unstructured, error) {
	port := viper.GetInt("notebook.defaultPort")
	return k8sClient.CreateVirtualService(notebook.Namespace, notebook.Name, "*", username, int32(port))
}

// DeleteNotebook godoc
// @Summary Delete a Notebook
// @Description Delete a Notebook by its ID
// @Tags notebook
// @Accept json
// @Produce json
// @Param id path int true "Notebook ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notebook/{id} [delete]
func DeleteNotebook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Invalid id parameter",
		})
		return
	}
	// Get the complete Notebook information
	notebook, err := model.GetNotebookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Notebook not found",
		})
		return
	}
	// Create K8s client
	k8sClient, err := services.NewK8s("./services/config")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	if err := deleteNotebookResources(k8sClient, notebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete Kubernetes resources: " + err.Error(),
		})
		return
	}

	if err := notebook.Delete(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notebook deleted successfully",
	})
}

// GetNotebook godoc
// @Summary Get a Notebook
// @Description Get a Notebook by its ID
// @Tags notebook
// @Accept json
// @Produce json
// @Param id path int true "Notebook ID"
// @Success 200 {object} NotebookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /notebook/{id} [get]
func GetNotebook(c *gin.Context) {

	idStr := c.Param("id")
	common.SysLog("id: " + idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Invalid id parameter",
		})
		return
	}

	notebook, err := model.GetNotebookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Notebook not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    notebook,
	})
}

// ListNotebooks godoc
// @Summary List Notebooks
// @Description Get a paginated list of Notebooks
// @Tags notebook
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} NotebooksResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /notebook/get-all [get]
func ListNotebooks(c *gin.Context) {
	// Get current user information
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated1",
		})
		return
	}
	id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated2",
		})
		return
	}
	userId, ok := id.(int)
	if !ok {
		// Handle type assertion failure
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize := common.ItemsPerPage
	offset := (page - 1) * pageSize

	var notebooks []model.Notebook
	var total int64
	var err error

	// Get notebooks based on user role
	if role == 10 || role == 100 {
		notebooks, total, err = model.GetAllNotebooksPaginated(offset, pageSize)
	} else {
		notebooks, total, err = model.GetUserNotebooksPaginated(userId, offset, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"notebooks": notebooks,
			"total":     total,
			"page":      page,
			"pageSize":  pageSize,
		},
	})
}

// ResetNotebook godoc
// @Summary Reset a Notebook
// @Description Reset a Notebook by its ID
// @Tags notebook
// @Accept json
// @Produce json
// @Param id path int true "Notebook ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notebook/reset/{id} [post]
func ResetNotebook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid id parameter",
		})
		return
	}

	// Retrieve the existing notebook
	existingNotebook, err := model.GetNotebookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Notebook not found!",
		})
		return
	}

	// Create K8s client
	k8sClient, err := services.NewK8s("./services/config")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create K8s client: " + err.Error(),
		})
		return
	}

	// Delete existing Kubernetes resources
	if err := deleteNotebookResources(k8sClient, existingNotebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete existing resources: " + err.Error(),
		})
		return
	}

	// Delete the notebook from the database
	existingNotebook.Status = "Resetting"
	existingNotebook.UpdatedAt = time.Now()
	if err := existingNotebook.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update notebook in database: " + err.Error(),
		})
		return
	}
	// Create a new notebook with the same information
	newNotebook := *existingNotebook
	newNotebook.ID = 0 // Reset ID for new insertion
	newNotebook.CreatedAt = time.Now()
	newNotebook.UpdatedAt = time.Now()

	// Create new Kubernetes resources
	if err := createNotebookResources(k8sClient, &newNotebook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create new resources: " + err.Error(),
		})
		return
	}
}

type NotebookUpdateRequest struct {
	ResourceCPU    *string `json:"resource_cpu,omitempty"`
	ResourceMemory *string `json:"resource_memory,omitempty"`
	ResourceGPU    *int64  `json:"resource_gpu,omitempty"`
	ServicePort    *int32  `json:"service_port,omitempty"`
	MountPath      *string `json:"mount_path,omitempty"`
}

// UpdateNotebook godoc
// @Summary Update a Notebook
// @Description Update a Notebook by its ID
// @Tags notebook
// @Accept json
// @Produce json
// @Param id path int true "Notebook ID"
// @Param notebook body NotebookUpdateRequest true "Notebook update details"
// @Success 200 {object} NotebookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /notebook/{id} [put]
func UpdateNotebook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid id parameter"})
		return
	}

	var updateReq NotebookUpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	notebook, err := model.GetNotebookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Notebook not found"})
		return
	}

	// Update Notebook model
	updated := updateNotebookModel(notebook, &updateReq)

	if !updated {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "No changes to update"})
		return
	}

	// Create K8s client
	k8sClient, err := services.NewK8s("./services/config")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create K8s client"})
		return
	}

	// Update Kubernetes resources
	if err := updateK8sResources(k8sClient, notebook, &updateReq); err != nil {
		common.SysError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update Kubernetes resources"})
		return
	}

	// update notebook in database
	if err := notebook.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update notebook in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notebook updated successfully", "data": notebook})
}

func updateNotebookModel(notebook *model.Notebook, updateReq *NotebookUpdateRequest) bool {
	updated := false

	if updateReq.ResourceCPU != nil {
		notebook.ResourceCPU = *updateReq.ResourceCPU
		updated = true
	}
	if updateReq.ResourceMemory != nil {
		notebook.ResourceMemory = *updateReq.ResourceMemory
		updated = true
	}
	if updateReq.ResourceGPU != nil {
		notebook.ResourceGPU = *updateReq.ResourceGPU
		updated = true
	}
	// Add update logic for other fields

	return updated
}

func updateK8sResources(k8sClient *services.K8s, notebook *model.Notebook, updateReq *NotebookUpdateRequest) error {
	// ctx := context.Background()

	// update Pod
	pod, err := k8sClient.GetPod(notebook.Namespace, notebook.Name)
	if err != nil {
		return err
	}
	newPod := pod.DeepCopy()
	updatePodSpec(newPod, updateReq)
	err = k8sClient.DeletePod(notebook.Namespace, notebook.Name)
	if err != nil {
		return fmt.Errorf("failed to delete old pod: %v", err)
	}

	_, err = k8sClient.CreatePod(notebook.Namespace, newPod)
	if err != nil {
		return fmt.Errorf("failed to create new pod: %v", err)
	}

	// update Service
	if updateReq.ServicePort != nil {
		svc, err := k8sClient.GetService(notebook.Namespace, notebook.Name)
		if err != nil {
			return err
		}
		svc.Spec.Ports[0].Port = *updateReq.ServicePort
		if _, err := k8sClient.UpdateService(notebook.Namespace, svc); err != nil {
			return err
		}
	}

	return nil
}
func updatePodSpec(pod *corev1.Pod, updateReq *NotebookUpdateRequest) {
	container := &pod.Spec.Containers[0]

	if updateReq.ResourceCPU != nil {
		quantity := resource.MustParse(*updateReq.ResourceCPU)
		container.Resources.Limits[corev1.ResourceCPU] = quantity
		container.Resources.Requests[corev1.ResourceCPU] = quantity
	}
	if updateReq.ResourceMemory != nil {
		quantity := resource.MustParse(*updateReq.ResourceMemory)
		container.Resources.Limits[corev1.ResourceMemory] = quantity
		container.Resources.Requests[corev1.ResourceMemory] = quantity
	}
	if updateReq.ResourceGPU != nil {
		quantity := resource.NewQuantity(int64(*updateReq.ResourceGPU), resource.DecimalSI)
		container.Resources.Limits["nvidia.com/gpu"] = *quantity
		container.Resources.Requests["nvidia.com/gpu"] = *quantity
	}
}
func deleteNotebookResources(k8sClient *services.K8s, notebook *model.Notebook) error {
	ctx := context.Background()

	// Delete Pod
	err := k8sClient.DeletePod(notebook.Namespace, notebook.Name)
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Pod: %v", err)
	}

	// Delete Service
	err = k8sClient.DeleteService(notebook.Namespace, notebook.Name, nil)
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete Service: %v", err)
	}

	// Delete VirtualService
	err = k8sClient.DeleteVirtualService(ctx, notebook.Namespace, notebook.Name)
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete VirtualService: %v", err)
	}

	return nil
}

func createNotebookResources(k8sClient *services.K8s, notebook *model.Notebook) error {
	username := notebook.Name // Assuming the notebook name is the username

	labels := map[string]string{
		"app":      notebook.Name,
		"pod-type": "notebook",
		"user":     username,
	}

	// Create Pod
	_, err := createPodForNotebook(k8sClient, notebook, username, labels)
	if err != nil {
		return fmt.Errorf("failed to create Pod: %v", err)
	}

	// Create Service
	_, err = createServiceForNotebook(k8sClient, notebook, labels)
	if err != nil {
		return fmt.Errorf("failed to create Service: %v", err)
	}

	// Create VirtualService
	_, err = createVsForNotebook(k8sClient, notebook, username)
	if err != nil {
		return fmt.Errorf("failed to create VirtualService: %v", err)
	}

	return nil
}

// Swagger model definitions
type NotebookResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Notebook created successfully"`
	Data    model.Notebook `json:"data"`
}

type NotebooksResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:""`
	Data    NotebooksListData `json:"data"`
}

type NotebooksListData struct {
	Notebooks []model.Notebook `json:"notebooks"`
	Total     int64            `json:"total" example:"10"`
	Page      int              `json:"page" example:"1"`
	PageSize  int              `json:"pageSize" example:"20"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"An error occurred"`
}
