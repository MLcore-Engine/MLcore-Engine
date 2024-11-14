package services

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateTritonDeployment creates a Kubernetes Deployment for Triton Server
func (k *K8s) CreateTritonDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	ctx := context.Background()

	if namespace == "" {
		namespace = deployment.Namespace
		if namespace == "" {
			return nil, fmt.Errorf("namespace must be specified either in the deployment or as a parameter")
		}
	}

	// Try to delete existing deployment
	err := k.clientset.AppsV1().Deployments(namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{GracePeriodSeconds: new(int64)})
	if err != nil && !k8serrors.IsNotFound(err) {
		fmt.Printf("Error deleting existing deployment: %v\n", err)
	}

	// Create deployment
	createdDeployment, err := k.clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %v", err)
	}

	time.Sleep(time.Second) // wait 1 second

	return createdDeployment, nil
}

// CreateTritonService creates a Kubernetes Service for Triton Server with NodePort
func (k *K8s) CreateTritonService(namespace string, service *corev1.Service) (*corev1.Service, error) {
	ctx := context.Background()

	if namespace == "" {
		namespace = service.Namespace
		if namespace == "" {
			return nil, fmt.Errorf("namespace must be specified either in the service or as a parameter")
		}
	}

	// Try to delete existing service
	err := k.clientset.CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		fmt.Printf("Error deleting existing service: %v\n", err)
	}

	// Create service
	createdService, err := k.clientset.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	return createdService, nil
}

// DeleteDeployment deletes a Kubernetes Deployment
func (k *K8s) DeleteDeployment(namespace, name string) error {
	ctx := context.Background()
	return k.clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// DeleteService deletes a Kubernetes Service
func (k *K8s) DeleteService2(namespace, name string) error {
	ctx := context.Background()
	return k.clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
