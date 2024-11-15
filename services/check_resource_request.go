package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceRequest struct {
	CPU    resource.Quantity
	Memory resource.Quantity
	GPU    resource.Quantity
}

func (k *K8s) CheckClusterResource(request ResourceRequest) (bool, error) {

	nodes, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get nodes: %v", err)
	}

	for _, node := range nodes.Items {
		if k.canScheduleOnNode(&node, request) {
			return true, nil
		}
	}
	return false, nil

}

func (k *K8s) canScheduleOnNode(node *corev1.Node, request ResourceRequest) bool {

	allocatable := node.Status.Allocatable
	if allocatable.Cpu().Cmp(request.CPU) < 0 ||
		allocatable.Memory().Cmp(request.Memory) < 0 {
		return false
	}

	// 检查 GPU 资源
	if request.GPU.Cmp(resource.MustParse("0")) > 0 {
		if gpuQuantity, exists := allocatable["nvidia.com/gpu"]; exists {
			if (&gpuQuantity).Cmp(request.GPU) < 0 {
				return false
			}
		} else if !request.GPU.IsZero() {
			return false
		}
	}
	// 获取节点上运行的所有 Pod
	pods, err := k.clientset.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + node.Name,
	})
	if err != nil {
		// 处理错误，这里简单地返回 false
		return false
	}

	// 计算节点上已使用的资源
	var usedCPU, usedMemory, usedGPU resource.Quantity
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			usedCPU.Add(*container.Resources.Requests.Cpu())
			usedMemory.Add(*container.Resources.Requests.Memory())
			if gpuQuantity, exists := container.Resources.Requests["nvidia.com/gpu"]; exists {
				usedGPU.Add(gpuQuantity)
			}
		}
	}

	// 计算剩余可用资源
	availableCPU := allocatable.Cpu().DeepCopy()
	availableCPU.Sub(usedCPU)
	availableMemory := allocatable.Memory().DeepCopy()
	availableMemory.Sub(usedMemory)

	// 检查是否有足够的剩余 CPU 和内存资源
	if availableCPU.Cmp(request.CPU) < 0 || availableMemory.Cmp(request.Memory) < 0 {
		return false
	}

	// 只有当请求的 GPU 大于 0 时才检查 GPU 资源
	if request.GPU.Cmp(resource.MustParse("0")) > 0 {
		if gpuQuantity, exists := allocatable["nvidia.com/gpu"]; exists {
			availableGPU := gpuQuantity.DeepCopy()
			availableGPU.Sub(usedGPU)
			if availableGPU.Cmp(request.GPU) < 0 {
				return false
			}
		} else {
			// 如果请求 GPU 但节点没有 GPU 资源，返回 false
			return false
		}
	}

	// 如果所有检查都通过，返回 true
	return true
}
