package services

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8s) CreatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {

	ctx := context.Background()
	// if namespace is nil use pod.Namespace
	if namespace == "" {
		namespace = pod.Namespace
		if namespace == "" {
			return nil, fmt.Errorf("namespace must be specified either in the pod or as a parameter")
		}
	}
	// try to delete existing pod
	err := k.clientset.CoreV1().Pods(namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{GracePeriodSeconds: new(int64)})
	if err != nil && !k8serrors.IsNotFound(err) {
		fmt.Printf("Error deleting existing pod: %v\n", err)
	}

	// create pod
	createdPod, err := k.clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %v", err)
	}

	time.Sleep(time.Second) // wait 1 second

	return createdPod, nil
}
