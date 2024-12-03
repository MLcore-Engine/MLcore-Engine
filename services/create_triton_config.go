package services

import (
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
)

type TritonConfig struct {
	// Server Configuration
	ModelRepository          string `json:"model_repository"`
	StrictModelConfig        bool   `json:"strict_model_config"`
	AllowPollModelRepository bool   `json:"allow_poll_model_repository"`
	PollRepoSeconds          int    `json:"poll_repo_seconds"`

	// HTTP Configuration
	HttpPort        int  `json:"http_port"`
	HttpThreadCount int  `json:"http_thread_count"`
	AllowHttp       bool `json:"allow_http"`

	// gRPC Configuration
	GrpcPort                    int  `json:"grpc_port"`
	GrpcInferAllocationPoolSize int  `json:"grpc_infer_allocation_pool_size"`
	AllowGrpc                   bool `json:"allow_grpc"`

	// Metrics Configuration
	AllowMetrics      bool `json:"allow_metrics"`
	MetricsPort       int  `json:"metrics_port"`
	MetricsIntervalMs int  `json:"metrics_interval_ms"`

	// GPU Configuration
	GpuMemoryFraction             float64 `json:"gpu_memory_fraction"`
	MinSupportedComputeCapability float64 `json:"min_supported_compute_capability"`

	// Logging Configuration
	LogVerbose int  `json:"log_verbose"`
	LogInfo    bool `json:"log_info"`
	LogWarning bool `json:"log_warning"`
	LogError   bool `json:"log_error"`
}

// GetTritonDeployment returns the Deployment object for Triton Server without InitContainer and Volumes
func GetTritonDeployment(name, namespace, image string, replicas int32, labels string, cpu, memory, gpu int64, mountPath string, config TritonConfig) (*appsv1.Deployment, error) {
	// Parse labels string to map
	var labelsMap map[string]string
	if err := json.Unmarshal([]byte(labels), &labelsMap); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %v", err)
	}

	args := []string{
		fmt.Sprintf("--model-repository=%s", config.ModelRepository),
		fmt.Sprintf("--strict-model-config=%t", config.StrictModelConfig),
	}

	if config.AllowPollModelRepository {
		args = append(args, "--allow-poll-model-repository=true")
		args = append(args, fmt.Sprintf("--poll-repo-seconds=%d", config.PollRepoSeconds))
	}

	// HTTP Configuration
	if config.AllowHttp {
		args = append(args, "--allow-http=true")
		args = append(args, fmt.Sprintf("--http-port=%d", config.HttpPort))
		args = append(args, fmt.Sprintf("--http-thread-count=%d", config.HttpThreadCount))
	}

	// gRPC Configuration
	if config.AllowGrpc {
		args = append(args, "--allow-grpc=true")
		args = append(args, fmt.Sprintf("--grpc-port=%d", config.GrpcPort))
		args = append(args, fmt.Sprintf("--grpc-infer-allocation-pool-size=%d", config.GrpcInferAllocationPoolSize))
	}

	// Metrics Configuration
	if config.AllowMetrics {
		args = append(args, "--allow-metrics=true")
		args = append(args, fmt.Sprintf("--metrics-port=%d", config.MetricsPort))
		args = append(args, fmt.Sprintf("--metrics-interval-ms=%d", config.MetricsIntervalMs))
	}

	// GPU Configuration
	if gpu > 0 {
		args = append(args, fmt.Sprintf("--gpu-memory-fraction=%.2f", config.GpuMemoryFraction))
		args = append(args, fmt.Sprintf("--min-supported-compute-capability=%.1f", config.MinSupportedComputeCapability))
	}

	// Logging Configuration
	args = append(args, fmt.Sprintf("--log-verbose=%d", config.LogVerbose))
	if config.LogInfo {
		args = append(args, "--log-info=true")
	}
	if config.LogWarning {
		args = append(args, "--log-warning=true")
	}
	if config.LogError {
		args = append(args, "--log-error=true")
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labelsMap,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labelsMap,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelsMap,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http-triton",
									ContainerPort: 8000,
								},
								{
									Name:          "grpc-triton",
									ContainerPort: 8001,
								},
								{
									Name:          "metrics-triton",
									ContainerPort: 8002,
								},
							},
							Command: []string{"tritonserver"},
							Args:    args,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resourceQuantity(cpu),
									corev1.ResourceMemory: resourceQuantity(memory * 1024 * 1024 * 1024), // Convert to bytes
									"nvidia.com/gpu":      resourceQuantity(gpu),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resourceQuantity(cpu),
									corev1.ResourceMemory: resourceQuantity(memory * 1024 * 1024 * 1024), // Convert to bytes
									"nvidia.com/gpu":      resourceQuantity(gpu),
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

// GetTritonService returns the Service object for Triton Server with NodePort
func GetTritonService(name, namespace string, labels string) (*corev1.Service, error) {
	// Parse labels string to map
	var labelsMap map[string]string
	if err := json.Unmarshal([]byte(labels), &labelsMap); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %v", err)
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labelsMap,
		},
		Spec: corev1.ServiceSpec{
			Selector: labelsMap,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8000,
					Name:       "http-triton",
					TargetPort: intstrFromInt(8000),
				},
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8001,
					Name:       "grpc-triton",
					TargetPort: intstrFromInt(8001),
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}, nil
}

// Helper function to create resource.Quantity
func resourceQuantity(amount int64) resource.Quantity {
	return *resource.NewQuantity(amount, resource.DecimalSI)
}

// Helper function to create intstr.IntOrString from int
func intstrFromInt(i int32) intstr.IntOrString {
	return intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: i,
	}
}
