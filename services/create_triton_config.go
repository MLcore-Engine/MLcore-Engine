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

// GetTritonDeployment returns the Deployment object for Triton Server without InitContainer and Volumes
func GetTritonDeployment(name, namespace, image string, replicas int32, labels string, cpu, memory, gpu int64, mountPath string) (*appsv1.Deployment, error) {
	// Parse labels string to map
	var labelsMap map[string]string
	if err := json.Unmarshal([]byte(labels), &labelsMap); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %v", err)
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
							Args: []string{
								"--model-repository=/model",
								"--allow-gpu-metrics=false",
								"--strict-model-config=false",
							},
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
