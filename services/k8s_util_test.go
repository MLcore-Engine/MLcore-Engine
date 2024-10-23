package services // 确保这与你的主包名一致

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	dfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGetPods(t *testing.T) {
	// 创建一个模拟的 Kubernetes 客户端
	clientset := fake.NewSimpleClientset()

	// 创建一个 K8s 实例
	k8s := &K8s{clientset: clientset}

	// 创建测试用的 namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}

	// 创建测试用的 pods
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod1",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app": "test",
				},
			},
			Spec: corev1.PodSpec{
				NodeName: "node1",
				Containers: []corev1.Container{
					{
						Name: "container1",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("100m"),
								corev1.ResourceMemory: resource.MustParse("128Mi"),
								"nvidia.com/gpu":      resource.MustParse("1"),
							},
						},
					},
				},
				NodeSelector: map[string]string{
					"disktype": "ssd",
				},
			},
			Status: corev1.PodStatus{
				Phase:     corev1.PodRunning,
				HostIP:    "192.168.1.1",
				PodIP:     "10.0.0.1",
				StartTime: &metav1.Time{Time: time.Now()},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod2",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app": "test",
				},
			},
			Status: corev1.PodStatus{
				Phase:  corev1.PodPending,
				HostIP: "192.168.1.2",
				PodIP:  "10.0.0.2",
			},
		},
	}

	for _, pod := range pods {
		_, err := clientset.CoreV1().Pods("test-namespace").Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error creating pod: %v", err)
		}
	}

	// 创建测试用的 service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "test",
			},
		},
	}
	_, err = clientset.CoreV1().Services("test-namespace").Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating service: %v", err)
	}

	// 创建 Endpoints
	endpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP: "10.0.0.1",
						TargetRef: &corev1.ObjectReference{
							Kind: "Pod",
							Name: "pod1",
						},
					},
					{
						IP: "10.0.0.2",
						TargetRef: &corev1.ObjectReference{
							Kind: "Pod",
							Name: "pod2",
						},
					},
				},
			},
		},
	}
	_, err = clientset.CoreV1().Endpoints("test-namespace").Create(context.TODO(), endpoints, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating endpoints: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name        string
		namespace   string
		serviceName string
		podName     string
		labels      map[string]string
		expectedLen int
	}{
		{"List all pods in namespace", "test-namespace", "", "", nil, 2},
		{"Get specific pod", "test-namespace", "", "pod1", nil, 1},
		{"Get pods for service", "test-namespace", "test-service", "", nil, 2},
		{"Get pods with label", "test-namespace", "", "", map[string]string{"app": "test"}, 2},
		{"No pods in non-existent namespace", "non-existent", "", "", nil, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pods, err := k8s.GetPods(tc.namespace, tc.serviceName, tc.podName, tc.labels)
			if err != nil {
				t.Fatalf("Error getting pods: %v", err)
			}
			if len(pods) != tc.expectedLen {
				t.Errorf("Expected %d pods, got %d", tc.expectedLen, len(pods))
			}

			// Additional checks for the first test case (all pods)
			if tc.name == "List all pods in namespace" && len(pods) > 0 {
				pod := pods[0]
				if pod.Name != "pod1" {
					t.Errorf("Expected pod name 'pod1', got '%s'", pod.Name)
				}
				if pod.HostIP != "192.168.1.1" {
					t.Errorf("Expected host IP '192.168.1.1', got '%s'", pod.HostIP)
				}
				if pod.PodIP != "10.0.0.1" {
					t.Errorf("Expected pod IP '10.0.0.1', got '%s'", pod.PodIP)
				}
				if pod.Status != "Running" {
					t.Errorf("Expected status 'Running', got '%s'", pod.Status)
				}
				if pod.NodeName != "node1" {
					t.Errorf("Expected node name 'node1', got '%s'", pod.NodeName)
				}
				if !reflect.DeepEqual(pod.Labels, map[string]string{"app": "test"}) {
					t.Errorf("Expected labels %v, got %v", map[string]string{"app": "test"}, pod.Labels)
				}
				if pod.Memory != 0.125 { // 128Mi = 0.125 Gi
					t.Errorf("Expected memory 0.125 Gi, got %f Gi", pod.Memory)
				}
				if pod.CPU != 0.1 { // 100m = 0.1 cores
					t.Errorf("Expected CPU 0.1 cores, got %f cores", pod.CPU)
				}
				if pod.GPU != 1 {
					t.Errorf("Expected GPU 1, got %d", pod.GPU)
				}
				if !reflect.DeepEqual(pod.NodeSelector, map[string]string{"disktype": "ssd"}) {
					t.Errorf("Expected node selector %v, got %v", map[string]string{"disktype": "ssd"}, pod.NodeSelector)
				}
			}
		})
	}
}

func TestGetPodIP(t *testing.T) {
	// 创建一个模拟的 Kubernetes 客户端
	clientset := fake.NewSimpleClientset()

	// 创建一个 K8s 实例
	k8s := &K8s{clientset: clientset}

	// 创建测试用的 namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cloudai-2",
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}

	// 创建测试用的 service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "face-search-vip-service",
			Namespace: "cloudai-2",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "face-search",
			},
		},
	}
	_, err = clientset.CoreV1().Services("cloudai-2").Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating service: %v", err)
	}

	// 创建测试用的 pods
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod1",
				Namespace: "cloudai-2",
				Labels: map[string]string{
					"app": "face-search",
				},
			},
			Status: corev1.PodStatus{
				PodIP: "10.0.0.1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod2",
				Namespace: "cloudai-2",
				Labels: map[string]string{
					"app": "face-search",
				},
			},
			Status: corev1.PodStatus{
				PodIP: "10.0.0.2",
			},
		},
	}

	// 创建 Endpoint
	endpoints := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "face-search-vip-service",
			Namespace: "cloudai-2",
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP: "10.0.0.1",
						TargetRef: &corev1.ObjectReference{
							Kind: "Pod",
							Name: "pod1",
						},
					},
					{
						IP: "10.0.0.2",
						TargetRef: &corev1.ObjectReference{
							Kind: "Pod",
							Name: "pod2",
						},
					},
				},
			},
		},
	}
	_, err = clientset.CoreV1().Endpoints("cloudai-2").Create(context.TODO(), endpoints, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating endpoints: %v", err)
	}

	for _, pod := range pods {
		_, err := clientset.CoreV1().Pods("cloudai-2").Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error creating pod: %v", err)
		}
	}

	// 测试 GetPodIP 函数
	podIPs, err := k8s.GetPodIP("", "")
	if err != nil {
		t.Fatalf("Error getting pod IPs: %v", err)
	}

	// 验证结果
	expectedIPs := []string{"10.0.0.1", "10.0.0.2"}
	if !reflect.DeepEqual(podIPs, expectedIPs) {
		t.Errorf("Expected pod IPs %v, but got %v", expectedIPs, podIPs)
	}

	// 测试使用特定的 namespace 和 service 名称
	podIPs, err = k8s.GetPodIP("cloudai-2", "face-search-vip-service")
	if err != nil {
		t.Fatalf("Error getting pod IPs with specific namespace and service: %v", err)
	}

	if !reflect.DeepEqual(podIPs, expectedIPs) {
		t.Errorf("Expected pod IPs %v, but got %v", expectedIPs, podIPs)
	}

	// 测试使用不存在的 namespace
	podIPs, err = k8s.GetPodIP("non-existent", "")
	if err != nil {
		t.Fatalf("Error getting pod IPs with non-existent namespace: %v", err)
	}

	if len(podIPs) != 0 {
		t.Errorf("Expected no pod IPs for non-existent namespace, but got %v", podIPs)
	}
}

type testEnv struct {
	clientset *fake.Clientset
	k8s       *K8s
	counter   int64
}

// setupTestEnv 创建并返回一个测试环境
func setupTestEnv() *testEnv {
	clientset := fake.NewSimpleClientset()
	k8s := &K8s{clientset: clientset}

	return &testEnv{
		clientset: clientset,
		k8s:       k8s,
		counter:   0,
	}
}

// createNamespace 创建一个新的命名空间并返回其名称
func (env *testEnv) createNamespace(t *testing.T) string {
	counter := atomic.AddInt64(&env.counter, 1)
	namespace := fmt.Sprintf("test-namespace-%d", counter)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := env.clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating namespace: %v", err)
	}
	return namespace
}

// createPods 在指定的命名空间中创建测试用的 pods
func (env *testEnv) createPods(t *testing.T, namespace string) []corev1.Pod {
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod1",
				Namespace: namespace,
				Labels: map[string]string{
					"app": "test",
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod2",
				Namespace: namespace,
				Labels: map[string]string{
					"app": "test",
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPending,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pod3",
				Namespace: namespace,
				Labels: map[string]string{
					"app": "other",
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
			},
		},
	}

	for _, pod := range pods {
		_, err := env.clientset.CoreV1().Pods(namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error creating pod: %v", err)
		}
	}

	return pods
}

func TestDeletePods(t *testing.T) {
	env := setupTestEnv()

	testCases := []struct {
		name        string
		serviceName string
		podName     string
		status      string
		labels      map[string]string
		expectedLen int
	}{
		{"Delete all pods", "", "", "", nil, 3},
		{"Delete pods by status", "", "", "Running", nil, 2},
		{"Delete pods by label", "", "", "", map[string]string{"app": "test"}, 2},
		{"Delete specific pod", "", "pod1", "", nil, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			namespace := env.createNamespace(t)
			env.createPods(t, namespace)

			deletedPods, err := env.k8s.DeletePods(namespace, tc.serviceName, tc.podName, tc.status, tc.labels)
			if err != nil {
				t.Fatalf("Error deleting pods: %v", err)
			}

			if len(deletedPods) != tc.expectedLen {
				t.Errorf("Expected to delete %d pods, but deleted %d", tc.expectedLen, len(deletedPods))
			}

			remainingPods, err := env.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				t.Fatalf("Error listing remaining pods: %v", err)
			}

			expectedRemaining := 3 - tc.expectedLen
			if len(remainingPods.Items) != expectedRemaining {
				t.Errorf("Expected %d remaining pods, but got %d", expectedRemaining, len(remainingPods.Items))
			}
		})
	}

	// 测试不存在的命名空间
	t.Run("Delete from non-existent namespace", func(t *testing.T) {
		deletedPods, err := env.k8s.DeletePods("non-existent", "", "", "", nil)
		if err != nil {
			t.Fatalf("Error deleting pods from non-existent namespace: %v", err)
		}
		if len(deletedPods) != 0 {
			t.Errorf("Expected 0 deleted pods from non-existent namespace, but got %d", len(deletedPods))
		}
	})
}

func TestGetNode(t *testing.T) {
	// 创建 fake clientset
	clientset := fake.NewSimpleClientset()

	// 创建 K8s 实例
	k8s := &K8s{clientset: clientset}

	// 创建测试节点
	node1 := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
			Labels: map[string]string{
				"env": "test",
			},
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "192.168.1.1"},
			},
			Allocatable: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("4"),
				corev1.ResourceMemory: resource.MustParse("8Gi"),
				"nvidia.com/gpu":      resource.MustParse("2"),
			},
		},
	}

	node2 := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node2",
			Labels: map[string]string{
				"env": "prod",
			},
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeInternalIP, Address: "192.168.1.2"},
			},
			Allocatable: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("8"),
				corev1.ResourceMemory: resource.MustParse("16Gi"),
				"nvidia.com/gpu":      resource.MustParse("4"),
			},
		},
	}

	// 将测试节点添加到 fake clientset
	_, err := clientset.CoreV1().Nodes().Create(context.TODO(), node1, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test node1: %v", err)
	}
	_, err = clientset.CoreV1().Nodes().Create(context.TODO(), node2, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test node2: %v", err)
	}
	// fmt.Println()
	// 定义测试用例
	tests := []struct {
		name          string
		label         string
		nodeName      string
		ip            string
		expectedNodes []NodeInfo
	}{
		{
			name:     "Get all nodes",
			label:    "",
			nodeName: "",
			ip:       "",
			expectedNodes: []NodeInfo{
				{CPU: 4, Memory: 8, GPU: 2, Labels: map[string]string{"env": "test"}, Name: "node1", HostIP: "192.168.1.1"},
				{CPU: 8, Memory: 16, GPU: 4, Labels: map[string]string{"env": "prod"}, Name: "node2", HostIP: "192.168.1.2"},
			},
		},
		{
			name:     "Get node by label",
			label:    "env=test",
			nodeName: "",
			ip:       "",
			expectedNodes: []NodeInfo{
				{CPU: 4, Memory: 8, GPU: 2, Labels: map[string]string{"env": "test"}, Name: "node1", HostIP: "192.168.1.1"},
			},
		},
		{
			name:     "Get node by name",
			label:    "",
			nodeName: "node2",
			ip:       "",
			expectedNodes: []NodeInfo{
				{CPU: 8, Memory: 16, GPU: 4, Labels: map[string]string{"env": "prod"}, Name: "node2", HostIP: "192.168.1.2"},
			},
		},
		{
			name:     "Get node by IP",
			label:    "",
			nodeName: "",
			ip:       "192.168.1.1",
			expectedNodes: []NodeInfo{
				{CPU: 4, Memory: 8, GPU: 2, Labels: map[string]string{"env": "test"}, Name: "node1", HostIP: "192.168.1.1"},
			},
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := k8s.GetNode(tt.label, tt.nodeName, tt.ip)
			if err != nil {
				t.Fatalf("GetNode() error = %v", err)
			}
			if !reflect.DeepEqual(nodes, tt.expectedNodes) {
				t.Errorf("GetNode() got = %v, want %v", nodes, tt.expectedNodes)
			}
		})
	}
}

func TestLabelNode(t *testing.T) {
	// 创建 fake clientset
	clientset := fake.NewSimpleClientset()

	// 创建 K8s 实例
	k8s := &K8s{clientset: clientset}

	// 创建测试节点
	node1 := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
			Labels: map[string]string{
				"existing": "label",
			},
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeHostName, Address: "node1"},
				{Type: corev1.NodeInternalIP, Address: "192.168.1.1"},
			},
		},
	}

	node2 := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node2",
		},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeHostName, Address: "node2"},
				{Type: corev1.NodeInternalIP, Address: "192.168.1.2"},
			},
		},
	}

	// 将测试节点添加到 fake clientset
	_, err := clientset.CoreV1().Nodes().Create(context.TODO(), node1, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test node1: %v", err)
	}
	_, err = clientset.CoreV1().Nodes().Create(context.TODO(), node2, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test node2: %v", err)
	}

	// 定义测试用例
	testCases := []struct {
		name          string
		ips           []string
		labels        map[string]string
		expectedIPs   []string
		expectedError bool
	}{
		{
			name:        "Label single node",
			ips:         []string{"192.168.1.1"},
			labels:      map[string]string{"env": "prod"},
			expectedIPs: []string{"192.168.1.1"},
		},
		{
			name:        "Label multiple nodes",
			ips:         []string{"192.168.1.1", "192.168.1.2"},
			labels:      map[string]string{"role": "worker"},
			expectedIPs: []string{"192.168.1.1", "192.168.1.2"},
		},
		{
			name:        "Label non-existent node",
			ips:         []string{"192.168.1.3"},
			labels:      map[string]string{"test": "label"},
			expectedIPs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updatedIPs, err := k8s.LabelNode(tc.ips, tc.labels)

			if tc.expectedError && err == nil {
				t.Errorf("Expected error, but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(updatedIPs, tc.expectedIPs) {
				t.Errorf("Expected IPs %v, but got %v", tc.expectedIPs, updatedIPs)
			}
			// another skill
			// if len(updatedIPs) != len(tc.expectedIPs) {
			// 	t.Errorf("Expected %d IPs, but got %d", len(tc.expectedIPs), len(updatedIPs))
			// } else {
			// 	for i, ip := range updatedIPs {
			// 		if ip != tc.expectedIPs[i] {
			// 			t.Errorf("IP mismatch at index %d: expected %s, got %s", i, tc.expectedIPs[i], ip)
			// 		}
			// 	}
			// }

			// 验证节点标签是否正确更新
			for _, ip := range tc.expectedIPs {
				var nodeName string
				if ip == "192.168.1.1" {
					nodeName = "node1"
				} else if ip == "192.168.1.2" {
					nodeName = "node2"
				}

				updatedNode, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
				fmt.Println(updatedNode.Labels)
				if err != nil {
					t.Errorf("Error getting updated node: %v", err)
					continue
				}

				for key, value := range tc.labels {
					if updatedNode.Labels[key] != value {
						t.Errorf("Expected label %s=%s on node %s, but got %s", key, value, nodeName, updatedNode.Labels[key])
					}
				}
			}
		})
	}
}

func TestDeleteService(t *testing.T) {
	// 创建一个 fake clientset
	clientset := fake.NewSimpleClientset()

	// 创建一个 K8s 实例，使用 fake clientset
	k := &K8s{
		clientset: clientset,
	}

	// 创建一个测试命名空间
	namespace := "test-namespace"

	// 创建一些测试服务
	service1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service1",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "test",
			},
		},
	}
	service2 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service2",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "test",
			},
		},
	}

	// 将测试服务添加到 fake clientset
	_, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), service1, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test service1: %v", err)
	}
	_, err = clientset.CoreV1().Services(namespace).Create(context.TODO(), service2, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test service2: %v", err)
	}

	// 测试删除指定名称的服务
	err = k.DeleteService(namespace, "service1", nil)
	if err != nil {
		t.Errorf("Error deleting service by name: %v", err)
	}

	// 验证 service1 已被删除
	_, err = clientset.CoreV1().Services(namespace).Get(context.TODO(), "service1", metav1.GetOptions{})
	if err == nil {
		t.Errorf("Service1 should have been deleted")
	}

	// 测试使用标签删除服务
	labels := map[string]string{"app": "test"}
	err = k.DeleteService(namespace, "", labels)
	if err != nil {
		t.Errorf("Error deleting services by labels: %v", err)
	}

	// 验证所有带有标签的服务都已被删除
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=test",
	})
	if err != nil {
		t.Errorf("Error listing services: %v", err)
	}
	if len(services.Items) != 0 {
		t.Errorf("Expected 0 services, but got %d", len(services.Items))
	}
}

func TestGetCRD(t *testing.T) {
	// 创建一个模拟的动态客户端
	scheme := runtime.NewScheme()
	gvr := schema.GroupVersionResource{Group: "test.com", Version: "v1", Resource: "testresources"}
	dynamicClient := dfake.NewSimpleDynamicClientWithCustomListKinds(scheme, map[schema.GroupVersionResource]string{
		gvr: "TestResourceList",
	})

	k8s := &K8s{
		dynamicClient: dynamicClient,
	}

	// 创建一个测试用的对象
	testObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "test.com/v1",
			"kind":       "TestResource",
			"metadata": map[string]interface{}{
				"name":      "test-resource",
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"replicas": int64(3), // 使用 int64 而不是 int
			},
			"status": map[string]interface{}{
				"phase": "Running",
			},
		},
	}

	// 将测试对象添加到模拟客户端
	_, err := dynamicClient.Resource(gvr).Namespace("default").Create(context.TODO(), testObj, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test object: %v", err)
	}

	// 测试 GetCRD 函数
	crdObj, err := k8s.GetCRD("test.com", "v1", "testresources", "default", "test-resource")
	if err != nil {
		t.Fatalf("Error getting CRD: %v", err)
	}

	if crdObj.Name != "test-resource" {
		t.Errorf("Expected name 'test-resource', got '%s'", crdObj.Name)
	}

	if crdObj.Status != "Running" {
		t.Errorf("Expected status 'Running', got '%s'", crdObj.Status)
	}
}

func TestGetVirtualService(t *testing.T) {
	// 假设配置文件名为 "config" 并位于当前目录
	kubeconfig := filepath.Join(".", "config")

	// 使用 kubeconfig 创建 config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating dynamic client: %v", err)
	}

	k8s := &K8s{
		dynamicClient: dynamicClient,
	}

	// 设置 VirtualService 的参数
	group := "networking.istio.io"
	version := "v1beta1"
	plural := "virtualservices"
	namespace := "jupyter"
	name := "notebook-jupyter-user0163-pod"

	// 获取 VirtualService
	vsObj, err := k8s.GetCRD(group, version, plural, namespace, name)
	if err != nil {
		t.Fatalf("Error getting VirtualService: %v", err)
	}

	if vsObj == nil {
		t.Fatalf("Received nil CRDObject")
	}

	// 验证获取到的 VirtualService
	if vsObj.Name != name {
		t.Errorf("Expected name '%s', got '%s'", name, vsObj.Name)
	}

	if vsObj.Namespace != namespace {
		t.Errorf("Expected namespace '%s', got '%s'", namespace, vsObj.Namespace)
	}

	// 输出获取到的 VirtualService 对象的详细信息，用于调试
	// t.Logf("Retrieved VirtualService object: %+v", vsObj)

	// 解析 Spec 字符串为 map
	var spec map[string]interface{}
	err = json.Unmarshal([]byte(vsObj.Spec), &spec)
	if err != nil {
		t.Fatalf("Failed to unmarshal Spec: %v", err)
	}

	// 验证 spec 中的字段
	gateways, ok := spec["gateways"].([]interface{})
	if !ok || len(gateways) == 0 {
		t.Errorf("Expected gateways to be non-empty slice")
	} else if gateways[0] != "kubeflow/kubeflow-gateway" {
		t.Errorf("Expected gateway to be 'kubeflow/kubeflow-gateway', got '%v'", gateways[0])
	}

	// 验证 http 路由规则
	http, ok := spec["http"].([]interface{})
	if !ok || len(http) == 0 {
		t.Errorf("Expected http to be non-empty slice")
	} else {
		httpRule, ok := http[0].(map[string]interface{})
		if !ok {
			t.Errorf("Expected http[0] to be a map")
		} else {
			timeout, ok := httpRule["timeout"].(string)
			if !ok || timeout != "300s" {
				t.Errorf("Expected timeout to be '300s', got '%v'", timeout)
			}
		}
	}
}
func TestDeleteCRD(t *testing.T) {
	// 创建一个模拟的动态客户端
	scheme := runtime.NewScheme()
	dynamicClient := dfake.NewSimpleDynamicClient(scheme)

	k8s := &K8s{
		dynamicClient: dynamicClient,
	}

	// 创建一个测试用的 unstructured 对象
	testObj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "test.com/v1",
			"kind":       "TestResource",
			"metadata": map[string]interface{}{
				"name":      "test-resource",
				"namespace": "default",
			},
		},
	}

	// 将测试对象添加到模拟客户端
	gvr := schema.GroupVersionResource{Group: "test.com", Version: "v1", Resource: "testresources"}
	_, err := dynamicClient.Resource(gvr).Namespace("default").Create(context.TODO(), testObj, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating test object: %v", err)
	}

	// 测试 DeleteCRD 函数
	err = k8s.DeleteCRD("test.com", "v1", "testresources", "default", "test-resource")
	if err != nil {
		t.Fatalf("Error deleting CRD: %v", err)
	}

	// 验证对象是否已被删除
	_, err = dynamicClient.Resource(gvr).Namespace("default").Get(context.TODO(), "test-resource", metav1.GetOptions{})
	if err == nil {
		t.Errorf("Expected error when getting deleted resource, but got nil")
	}
}

func TestGetCRDStatus(t *testing.T) {
	k8s := &K8s{}

	testCases := []struct {
		name     string
		obj      *unstructured.Unstructured
		group    string
		plural   string
		expected string
	}{
		{
			name: "Workflow",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"phase": "Running",
						"nodes": map[string]interface{}{
							"node1": map[string]interface{}{
								"phase": "Succeeded",
							},
						},
					},
				},
			},
			group:    "argoproj.io",
			plural:   "workflows",
			expected: "Running",
		},
		{
			name: "Notebook",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"status": map[string]interface{}{
						"conditions": []interface{}{
							map[string]interface{}{
								"type": "Ready",
							},
						},
					},
				},
			},
			group:    "kubeflow.org",
			plural:   "notebooks",
			expected: "Ready",
		},
		// 可以添加更多测试用例，如 InferenceService 和 VolcanoJob
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, err := k8s.GetCRDStatus(tc.obj, tc.group, tc.plural)
			if err != nil {
				t.Fatalf("Error getting CRD status: %v", err)
			}
			if status != tc.expected {
				t.Errorf("Expected status '%s', got '%s'", tc.expected, status)
			}
		})
	}
}

func TestDeleteExistingVirtualService(t *testing.T) {
	// 假设配置文件名为 "config" 并位于当前目录
	kubeconfig := filepath.Join(".", "config")

	// 使用 kubeconfig 创建 config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating dynamic client: %v", err)
	}

	k8s := &K8s{
		dynamicClient: dynamicClient,
	}

	// 设置 VirtualService 的参数
	group := "networking.istio.io"
	version := "v1beta1"
	plural := "virtualservices"
	namespace := "jupyter"
	name := "notebook-jupyter-luowei234-pod"

	// 测试 DeleteCRD 函数
	err = k8s.DeleteCRD(group, version, plural, namespace, name)
	if err != nil {
		t.Fatalf("Error deleting VirtualService: %v", err)
	}

	// 验证 VirtualService 已被删除
	gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: plural}
	_, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err == nil {
		t.Errorf("Expected error when getting deleted VirtualService, but got nil")
	} else if !k8serrors.IsNotFound(err) {
		t.Errorf("Expected 'not found' error, but got: %v", err)
	}

	// 再次尝试删除同一 VirtualService，应该返回错误
	err = k8s.DeleteCRD(group, version, plural, namespace, name)
	if err == nil {
		t.Errorf("Expected error when deleting non-existent VirtualService, but got nil")
	} else if !k8serrors.IsNotFound(err) {
		t.Errorf("Expected 'not found' error, but got: %v", err)
	}
}

func TestCreateService(t *testing.T) {
	// 使用当前目录下的 kubeconfig 文件
	kubeconfig := filepath.Join(".", "config")

	// 使用 kubeconfig 创建 config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	k8s := &K8s{
		clientset: clientset,
	}

	// 设置 Service 的参数
	namespace := "jupyter"
	name := "luowei234-pod" // 使用测试前缀以避免冲突
	username := "luowei234"
	ports := []interface{}{3000}
	selector := map[string]string{
		"app":      "luowei234-pod",
		"pod-type": "notebook",
		"user":     "luowei234",
	}
	serviceType := corev1.ServiceTypeClusterIP
	var externalIP []string
	annotations := map[string]string{}
	var loadBalancerIP string
	disableLoadBalancer := false

	// 调用 CreateService 函数
	createdService, err := k8s.CreateService(
		namespace, name, username, ports, selector, serviceType,
		externalIP, annotations, loadBalancerIP, disableLoadBalancer,
	)

	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// 验证创建的 Service
	if createdService.Name != name {
		t.Errorf("Expected service name %s, got %s", name, createdService.Name)
	}
	if createdService.Namespace != namespace {
		t.Errorf("Expected namespace %s, got %s", namespace, createdService.Namespace)
	}
	// fmt.Println("-------------------------####-------------------------")
	// fmt.Println(createdService)
	// fmt.Println("test到这里了")
	// 验证 labels
	expectedLabels := map[string]string{
		"app":      "luowei234-pod",
		"pod-type": "notebook",
		"user":     "luowei234",
	}
	if !reflect.DeepEqual(createdService.Labels, expectedLabels) {
		t.Errorf("Labels do not match. Expected %v, got %v", expectedLabels, createdService.Labels)
	}

	// 验证 spec
	if createdService.Spec.Type != corev1.ServiceTypeClusterIP {
		t.Errorf("Expected service type ClusterIP, got %v", createdService.Spec.Type)
	}

	if len(createdService.Spec.Ports) != 1 {
		t.Fatalf("Expected 1 port, got %d", len(createdService.Spec.Ports))
	}

	expectedPort := corev1.ServicePort{
		Name:       "http0",
		Port:       3000,
		Protocol:   corev1.ProtocolTCP,
		TargetPort: intstr.FromInt(3000),
	}
	if !reflect.DeepEqual(createdService.Spec.Ports[0], expectedPort) {
		t.Errorf("Port does not match. Expected %v, got %v", expectedPort, createdService.Spec.Ports[0])
	}

	if !reflect.DeepEqual(createdService.Spec.Selector, selector) {
		t.Errorf("Selector does not match. Expected %v, got %v", selector, createdService.Spec.Selector)
	}

	// 验证 ClusterIP 不是 "None"
	if createdService.Spec.ClusterIP == "None" {
		t.Errorf("Expected ClusterIP to not be 'None', got 'None'")
	}

	// 从 API 服务器获取 Service 并进行额外验证
	retrievedService, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to retrieve service: %v", err)
	}

	if !reflect.DeepEqual(createdService, retrievedService) {
		t.Errorf("Retrieved service does not match created service")
	}

	// 清理：删除创建的 Service
	// err = clientset.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	// if err != nil {
	// 	t.Fatalf("Failed to delete test service: %v", err)
	// }
}
