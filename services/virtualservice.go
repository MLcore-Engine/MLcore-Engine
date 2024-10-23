package services

import (
	"context"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (k *K8s) CreateVirtualService(namespace, crdName, host, username string, port int32) (*unstructured.Unstructured, error) {

	virtualServiceGVR := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1beta1",
		Resource: "virtualservices",
	}

	vs := map[string]interface{}{
		"apiVersion": "networking.istio.io/v1beta1",
		"kind":       "VirtualService",
		"metadata": map[string]interface{}{
			"name":      crdName,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"gateways": []string{"kubeflow/kubeflow-gateway"},
			"hosts":    []string{host},
			"http": []map[string]interface{}{
				{
					"match": []map[string]interface{}{
						{
							"uri": map[string]interface{}{
								"prefix": fmt.Sprintf("/notebook/%s/%s/", namespace, username),
							},
						},
					},
					"rewrite": map[string]interface{}{
						"uri": fmt.Sprintf("/notebook/%s/%s/", namespace, username),
					},
					"route": []map[string]interface{}{
						{
							"destination": map[string]interface{}{
								"host": fmt.Sprintf("%s-svc.%s.svc.cluster.local", crdName, namespace),
								"port": map[string]interface{}{
									"number": port, // 直接使用 int32 类型，不要转换为字符串
								},
							},
						},
					},
					"timeout": "300s",
				},
			},
		},
	}

	existingVS, err := k.dynamicClient.Resource(virtualServiceGVR).Namespace(namespace).Get(context.TODO(), crdName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// VirtualService 不存在，创建新的
			createdVS, err := k.dynamicClient.Resource(virtualServiceGVR).Namespace(namespace).Create(context.TODO(), &unstructured.Unstructured{Object: vs}, metav1.CreateOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to create VirtualService: %v", err)
			}
			return createdVS, nil
		}
		return nil, fmt.Errorf("failed to get existing VirtualService: %v", err)
	}

	existingVS.Object["spec"] = vs["spec"]
	updatedVS, err := k.dynamicClient.Resource(virtualServiceGVR).Namespace(namespace).Update(context.TODO(), existingVS, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update VirtualService: %v", err)
	}

	return updatedVS, nil
}

func (k *K8s) DeleteVirtualService(ctx context.Context, namespace, name string) error {
	// Define the GVR for VirtualService
	virtualServiceGVR := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "virtualservices",
	}

	// Delete the VirtualService
	err := k.dynamicClient.Resource(virtualServiceGVR).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete VirtualService %s in namespace %s: %v", name, namespace, err)
	}

	return nil
}
