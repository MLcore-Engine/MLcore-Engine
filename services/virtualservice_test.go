package services

import (
	"path/filepath"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateVirtualService(t *testing.T) {
	// 创建一个模拟的动态客户端
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
	namespace := "jupyter"
	crdName := "notebook-jupyter-luowei234-pod"
	host := "*"
	notebookName := "luowei234"
	port := int32(3000)

	// 调用 CreateVirtualService 函数
	createdVS, err := k8s.CreateVirtualService(namespace, crdName, host, notebookName, port)
	if err != nil {
		t.Fatalf("Failed to create VirtualService: %v", err)
	}

	// 打印创建的对象，用于调试
	t.Logf("Created VirtualService: %+v", createdVS)

	// 验证创建的 VirtualService
	if createdVS.GetName() != crdName {
		t.Errorf("Expected name %s, got %s", crdName, createdVS.GetName())
	}
	if createdVS.GetNamespace() != namespace {
		t.Errorf("Expected namespace %s, got %s", namespace, createdVS.GetNamespace())
	}

	// 验证 spec 字段
	spec, found, err := unstructured.NestedMap(createdVS.Object, "spec")
	if !found || err != nil {
		t.Fatalf("Failed to get spec: found=%v, err=%v", found, err)
	}

	// 验证 gateways
	gateways, found, err := unstructured.NestedStringSlice(spec, "gateways")
	if !found || err != nil {
		t.Fatalf("Failed to get gateways: found=%v, err=%v", found, err)
	}
	if len(gateways) != 1 || gateways[0] != "kubeflow/kubeflow-gateway" {
		t.Errorf("Expected gateways [kubeflow/kubeflow-gateway], got %v", gateways)
	}

	// 验证 hosts
	hosts, found, err := unstructured.NestedStringSlice(spec, "hosts")
	if !found || err != nil {
		t.Fatalf("Failed to get hosts: found=%v, err=%v", found, err)
	}
	if len(hosts) != 1 || hosts[0] != "*" {
		t.Errorf("Expected hosts [*], got %v", hosts)
	}

	// 验证 http 配置
	http, found, err := unstructured.NestedSlice(spec, "http")
	if !found || err != nil {
		t.Fatalf("Failed to get http: found=%v, err=%v", found, err)
	}
	if len(http) != 1 {
		t.Fatalf("Expected 1 http rule, got %d", len(http))
	}

	// 可以继续添加更多的验证逻辑来检查 http 规则的详细内容
}
