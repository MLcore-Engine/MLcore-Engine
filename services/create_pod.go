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

	// 如果没有指定命名空间，使用 Pod 定义中的命名空间
	if namespace == "" {
		namespace = pod.Namespace
		if namespace == "" {
			return nil, fmt.Errorf("namespace must be specified either in the pod or as a parameter")
		}
	}
	pod.ObjectMeta.ResourceVersion = ""
	pod.ObjectMeta.UID = ""
	pod.ObjectMeta.CreationTimestamp = metav1.Time{}
	// 尝试删除已存在的 Pod
	err := k.clientset.CoreV1().Pods(namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{GracePeriodSeconds: new(int64)})
	if err != nil && !k8serrors.IsNotFound(err) {
		fmt.Printf("Error deleting existing pod: %v\n", err)
	}

	// 确保 Pod 使用指定的命名空间
	pod.Namespace = namespace

	// 创建 Pod
	createdPod, err := k.clientset.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %v", err)
	}

	time.Sleep(time.Second) // 等待 1 秒

	return createdPod, nil
}
