package services

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCheckClusterResource(t *testing.T) {
	// 加载 Kubernetes 配置
	config, err := clientcmd.BuildConfigFromFlags("", "config")
	if err != nil {
		t.Fatalf("Failed to build config: %v", err)
	}

	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create clientset: %v", err)
	}

	// 创建 K8s 服务实例
	k8sService := &K8s{
		clientset: clientset,
	}

	// 定义测试用例
	testCases := []struct {
		name            string
		resourceRequest ResourceRequest
		expectedResult  bool
	}{
		{
			name: "Small resource request",
			resourceRequest: ResourceRequest{
				CPU:    resource.MustParse("100m"),
				Memory: resource.MustParse("100Mi"),
				GPU:    resource.MustParse("0"),
			},
			expectedResult: true,
		},
		{
			name: "Large resource request",
			resourceRequest: ResourceRequest{
				CPU:    resource.MustParse("100"),
				Memory: resource.MustParse("100Gi"),
				GPU:    resource.MustParse("4"),
			},
			expectedResult: false,
		},
		{
			name: "GPU request",
			resourceRequest: ResourceRequest{
				CPU:    resource.MustParse("1"),
				Memory: resource.MustParse("1Gi"),
				GPU:    resource.MustParse("1"),
			},
			expectedResult: false, // 假设集群没有 GPU 资源
		},
	}

	// 运行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := k8sService.CheckClusterResource(tc.resourceRequest)
			if err != nil {
				t.Fatalf("Error checking cluster resource: %v", err)
			}
			if result != tc.expectedResult {
				t.Errorf("Expected %v, but got %v", tc.expectedResult, result)
			}
		})
	}
}
