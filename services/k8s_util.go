package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8s struct {
	clientset     kubernetes.Interface
	dynamicClient dynamic.Interface
	config        *rest.Config
}

type PodInfo struct {
	Name         string                 `json:"name"`
	HostIP       string                 `json:"host_ip"`
	PodIP        string                 `json:"pod_ip"`
	Status       string                 `json:"status"`
	StatusMore   map[string]interface{} `json:"status_more"`
	NodeName     string                 `json:"node_name"`
	Labels       map[string]string      `json:"labels"`
	Memory       float64                `json:"memory"`
	CPU          float64                `json:"cpu"`
	GPU          int                    `json:"gpu"`
	StartTime    time.Time              `json:"start_time"`
	NodeSelector map[string]string      `json:"node_selector"`
}

type NodeInfo struct {
	CPU    int
	Memory int
	GPU    int
	Labels map[string]string
	Name   string
	HostIP string
}

type CRDObject struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Annotations string `json:"annotations"`
	Labels      string `json:"labels"`
	Spec        string `json:"spec"`
	CreateTime  string `json:"create_time"`
	FinishTime  string `json:"finish_time"`
	Status      string `json:"status"`
	StatusMore  string `json:"status_more"`
}

func NewK8s(filePath string) (*K8s, error) {
	var config *rest.Config
	var err error

	if filePath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", filePath)
	} else {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig != "" {
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		} else {
			home := homedir.HomeDir()
			kubeconfig = filepath.Join(home, ".kube", "config")
			if _, err := os.Stat(kubeconfig); err == nil {
				config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			} else {
				config, err = rest.InClusterConfig()
			}
		}
	}

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8s{clientset: clientset, config: config, dynamicClient: dynamicClient}, nil
}

func (k *K8s) GetPods(namespace, serviceName, podName string, labels map[string]string) ([]PodInfo, error) {
	var allPods []corev1.Pod
	ctx := context.Background()

	if namespace != "" && serviceName == "" && podName == "" && len(labels) == 0 {
		pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		allPods = pods.Items
	} else if namespace != "" && podName != "" {
		pod, err := k.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		allPods = append(allPods, *pod)
	} else if namespace != "" && serviceName != "" {
		endpoints, err := k.clientset.CoreV1().Endpoints(namespace).Get(ctx, serviceName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if len(endpoints.Subsets) > 0 {
			for _, address := range endpoints.Subsets[0].Addresses {
				if address.TargetRef != nil && address.TargetRef.Kind == "Pod" {
					pod, err := k.clientset.CoreV1().Pods(namespace).Get(ctx, address.TargetRef.Name, metav1.GetOptions{})
					if err != nil {
						return nil, err
					}
					allPods = append(allPods, *pod)
				}
			}
		}
	} else if namespace != "" && len(labels) > 0 {
		labelSelector := metav1.LabelSelector{MatchLabels: labels}
		pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(&labelSelector),
		})
		if err != nil {
			return nil, err
		}
		allPods = pods.Items
	}

	var backPods []PodInfo
	for _, pod := range allPods {
		podInfo := PodInfo{
			Name:     pod.Name,
			HostIP:   pod.Status.HostIP,
			PodIP:    pod.Status.PodIP,
			Status:   string(pod.Status.Phase),
			NodeName: pod.Spec.NodeName,
			Labels:   pod.Labels,
		}

		// Convert status to map
		statusMore, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod.Status)
		if err != nil {
			fmt.Printf("Error converting status: %v\n", err)
		} else {
			podInfo.StatusMore = statusMore
		}

		// Calculate resources
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				podInfo.Memory += parseMemory(container.Resources.Requests.Memory().String())
				podInfo.CPU += parseCPU(container.Resources.Requests.Cpu().String())
				podInfo.GPU += parseGPU(container.Resources.Requests["nvidia.com/gpu"])
			}
		}

		// Parse start time
		if pod.Status.StartTime != nil {
			podInfo.StartTime = pod.Status.StartTime.Add(8 * time.Hour).Local()
		}

		// Parse node selector
		podInfo.NodeSelector = make(map[string]string)
		if pod.Spec.NodeSelector != nil {
			for k, v := range pod.Spec.NodeSelector {
				podInfo.NodeSelector[k] = v
			}
		}
		if pod.Spec.Affinity != nil && pod.Spec.Affinity.NodeAffinity != nil {
			if required := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution; required != nil {
				for _, term := range required.NodeSelectorTerms {
					for _, expr := range term.MatchExpressions {
						if expr.Operator == corev1.NodeSelectorOpIn && len(expr.Values) > 0 {
							podInfo.NodeSelector[expr.Key] = expr.Values[0]
						}
					}
				}
			}
		}

		backPods = append(backPods, podInfo)
	}

	return backPods, nil
}

func parseMemory(mem string) float64 {
	mem = strings.TrimSpace(mem)
	if mem == "" {
		return 0
	}

	value, err := strconv.ParseFloat(mem[:len(mem)-2], 64)
	if err != nil {
		return 0
	}

	unit := strings.ToLower(mem[len(mem)-2:])
	switch unit {
	case "ki":
		return value / (1024 * 1024)
	case "mi":
		return value / 1024
	case "gi":
		return value
	case "ti":
		return value * 1024
	default:
		// 假设输入是直接的字节数
		return value / (1024 * 1024 * 1024)
	}
}

func parseCPU(cpu string) float64 {
	cpu = strings.TrimSpace(cpu)
	if cpu == "" {
		return 0
	}

	if strings.HasSuffix(cpu, "m") {
		value, err := strconv.ParseFloat(cpu[:len(cpu)-1], 64)
		if err != nil {
			return 0
		}
		return value / 1000
	}

	value, err := strconv.ParseFloat(cpu, 64)
	if err != nil {
		return 0
	}
	return value
}

func parseGPU(gpu resource.Quantity) int {
	return int(gpu.Value())
}

func (k *K8s) GetPodIP(namespace, serviceName string) ([]string, error) {
	if namespace == "" {
		namespace = "cloudai-2"
	}
	if serviceName == "" {
		serviceName = "face-search-vip-service"
	}

	pods, err := k.GetPods(namespace, serviceName, "", nil)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error getting pods: %v", err)
	}

	var podIPs []string
	for _, pod := range pods {
		if pod.PodIP != "" {
			podIPs = append(podIPs, pod.PodIP)
		}
	}

	return podIPs, nil
}

// DeletePod 删除指定命名空间中的 Pod
func (k *K8s) DeletePod(namespace, podName string) error {
	ctx := context.Background()
	deleteOptions := metav1.DeleteOptions{}

	err := k.clientset.CoreV1().Pods(namespace).Delete(ctx, podName, deleteOptions)
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete pod %s in namespace %s: %v", podName, namespace, err)
	}

	return nil
}

func (k *K8s) DeletePods(namespace, serviceName, podName, status string, labels map[string]string) ([]PodInfo, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	allPods, err := k.GetPods(namespace, serviceName, podName, labels)
	if err != nil {
		return nil, fmt.Errorf("error getting pods: %w", err)
	}

	if status != "" {
		var filteredPods []PodInfo
		for _, pod := range allPods {
			if pod.Status == status {
				filteredPods = append(filteredPods, pod)
			}
		}
		allPods = filteredPods
	}

	deletedPods := []PodInfo{}
	for _, pod := range allPods {
		gracePeriod := int64(0)
		deleteOptions := metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriod,
		}

		err := k.clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod.Name, deleteOptions)
		if err != nil {
			fmt.Printf("Error deleting pod %s: %v\n", pod.Name, err)
		} else {
			fmt.Printf("Deleted pod %s\n", pod.Name)
			deletedPods = append(deletedPods, pod)
		}
	}

	return deletedPods, nil
}

func (k *K8s) GetNode(label, name, ip string) ([]NodeInfo, error) {
	var backNodes []NodeInfo

	listOptions := metav1.ListOptions{}
	if label != "" {
		listOptions.LabelSelector = label
	}

	nodeList, err := k.clientset.CoreV1().Nodes().List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("error listing nodes: %v", err)
	}

	for _, node := range nodeList.Items {
		backNode := NodeInfo{
			Labels: node.Labels,
			Name:   node.Name,
		}

		// CPU
		// cpu := node.Status.Allocatable.Cpu().String()
		// if strings.HasSuffix(cpu, "m") {
		// 	cpuValue, _ := strconv.Atoi(strings.TrimSuffix(cpu, "m"))
		// 	backNode.CPU = cpuValue / 1000
		// } else {
		// 	backNode.CPU, _ = strconv.Atoi(cpu)
		// }
		cpuQuantity := node.Status.Allocatable[corev1.ResourceCPU]
		backNode.CPU = int(cpuQuantity.Value())

		// Memory
		memoryQuantity := node.Status.Allocatable[corev1.ResourceMemory]
		memoryBytes := memoryQuantity.Value()
		backNode.Memory = int(memoryBytes / (1024 * 1024 * 1024)) // 转换为 GB

		// GPU
		// gpuStr := node.Status.Allocatable.Name("nvidia.com/gpu", resource.DecimalSI)
		// backNode.GPU, _ = strconv.Atoi(gpuStr.String())

		gpuQuantity := node.Status.Allocatable["nvidia.com/gpu"]
		backNode.GPU = int(gpuQuantity.Value())

		// HostIP
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				backNode.HostIP = address.Address
				break
			}
		}

		if (name == "" && ip == "") ||
			(name != "" && backNode.Name == name) ||
			(ip != "" && backNode.HostIP == ip) {
			backNodes = append(backNodes, backNode)
		}
	}

	return backNodes, nil
}

func (k *K8s) LabelNode(ips []string, labels map[string]string) ([]string, error) {
	allNodeIPs := []string{}
	// var allNodeIPs []string
	// 列出所有节点
	nodeList, err := k.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing nodes: %v", err)
	}

	for _, node := range nodeList.Items {
		var hostname, internalIP string

		// 获取节点的 Hostname 和 InternalIP
		for _, address := range node.Status.Addresses {
			switch address.Type {
			case corev1.NodeHostName:
				hostname = address.Address
			case corev1.NodeInternalIP:
				internalIP = address.Address
			}
		}

		// 检查节点 IP 是否在目标 IP 列表中
		if containsIP(ips, internalIP) {
			// 准备更新标签的补丁
			patchBytes, err := prepareLabelsPatch(labels)
			if err != nil {
				return nil, fmt.Errorf("error preparing patch for node %s: %v", hostname, err)
			}

			// 应用补丁更新节点标签
			_, err = k.clientset.CoreV1().Nodes().Patch(context.TODO(), hostname, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
			if err != nil {
				return nil, fmt.Errorf("error patching node %s: %v", hostname, err)
			}

			allNodeIPs = append(allNodeIPs, internalIP)
		}
	}

	return allNodeIPs, nil
}

// containsIP 检查 IP 是否在列表中
func containsIP(ips []string, ip string) bool {
	for _, i := range ips {
		if i == ip {
			return true
		}
	}
	return false
}

// prepareLabelsPatch 准备用于更新标签的补丁
func prepareLabelsPatch(labels map[string]string) ([]byte, error) {
	patchBody := map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": labels,
		},
	}

	patchBytes, err := json.Marshal(patchBody)
	if err != nil {
		return nil, err
	}

	return patchBytes, nil
}

func (k *K8s) CreateServiceForNotebook(namespace, name string, port int, labels map[string]string) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeNodePort,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http0",
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	// try create Service
	createdService, err := k.clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			// if Service already exists, update it
			existingService, err := k.clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to get existing service: %v", err)
			}
			existingService.Spec = service.Spec
			updatedService, err := k.clientset.CoreV1().Services(namespace).Update(context.TODO(), existingService, metav1.UpdateOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to update service: %v", err)
			}
			return updatedService, nil
		}
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	return createdService, nil
}

// func (k *K8s) CreateService(
// 	namespace, name, username string,
// 	ports []interface{},
// 	selector map[string]string,
// 	serviceType corev1.ServiceType,
// 	externalIP []string,
// 	annotations map[string]string,
// 	loadBalancerIP string,
// 	disableLoadBalancer bool,
// ) (*corev1.Service, error) {
// 	servicePorts := []corev1.ServicePort{}
// 	for i, port := range ports {
// 		var servicePort corev1.ServicePort
// 		switch p := port.(type) {
// 		case []interface{}:
// 			if len(p) > 1 {
// 				portNum, _ := strconv.Atoi(fmt.Sprint(p[0]))
// 				targetPort, _ := strconv.Atoi(fmt.Sprint(p[1]))
// 				servicePort = corev1.ServicePort{
// 					Name:       fmt.Sprintf("http%d", i),
// 					Port:       int32(portNum),
// 					Protocol:   corev1.ProtocolTCP,
// 					TargetPort: intstr.FromInt(targetPort),
// 				}
// 				if serviceType == corev1.ServiceTypeNodePort {
// 					servicePort.NodePort = int32(portNum)
// 				}
// 			}
// 		default:
// 			portNum, _ := strconv.Atoi(fmt.Sprint(p))
// 			servicePort = corev1.ServicePort{
// 				Name:       fmt.Sprintf("http%d", i),
// 				Port:       int32(portNum),
// 				Protocol:   corev1.ProtocolTCP,
// 				TargetPort: intstr.FromInt(portNum),
// 			}
// 			if serviceType == corev1.ServiceTypeNodePort {
// 				servicePort.NodePort = int32(portNum)
// 			}
// 		}
// 		servicePorts = append(servicePorts, servicePort)
// 	}

// 	var clusterIP string
// 	if disableLoadBalancer {
// 		clusterIP = "None"
// 	}

// 	service := &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:        name,
// 			Namespace:   namespace,
// 			Labels:      selector,
// 			Annotations: annotations,
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Ports:          servicePorts,
// 			Selector:       selector,
// 			Type:           serviceType,
// 			ExternalIPs:    externalIP,
// 			LoadBalancerIP: loadBalancerIP,
// 			ClusterIP:      clusterIP,
// 		},
// 	}

// 	// 尝试读取现有的 Service
// 	fmt.Printf("Attempting to create/update service %s in namespace %s\n", name, namespace)

// 	existingService, err := k.clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
// 	fmt.Printf("Get existing service error: %v\n", err)
// 	fmt.Println(existingService)
// 	if err != nil {
// 		if k8serrors.IsNotFound(err) {
// 			fmt.Println("Service not found, creating new service")
// 			createdService, err := k.clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
// 			fmt.Printf("Create service result: %+v\n", createdService)
// 			fmt.Printf("Create service error: %v\n", err)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to create service: %v", err)
// 			}
// 			return createdService, nil
// 		}
// 		return nil, fmt.Errorf("failed to get service: %v", err)
// 	}

// 	fmt.Println("Service exists, updating")
// 	updatedService, err := k.clientset.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
// 	fmt.Printf("Update service result: %+v\n", updatedService)
// 	fmt.Printf("Update service error: %v\n", err)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update service: %v", err)
// 	}

// 	return updatedService, nil
// }

func (k *K8s) DeleteService(namespace string, name string, labels map[string]string) error {
	if name != "" {
		err := k.clientset.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
			GracePeriodSeconds: new(int64), // 0 seconds
		})
		if err != nil {
			if !k8serrors.IsNotFound(err) {
				return fmt.Errorf("error deleting service %s: %v", name, err)
			}
		} else {
			fmt.Printf("Deleted service: %s in namespace: %s\n", name, namespace)
		}
	}

	if len(labels) > 0 {
		labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: labels,
		})

		services, err := k.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			return fmt.Errorf("error listing services: %v", err)
		}

		for _, service := range services.Items {
			err := k.clientset.CoreV1().Services(service.Namespace).Delete(context.TODO(), service.Name, metav1.DeleteOptions{
				GracePeriodSeconds: new(int64), // 0 seconds
			})
			if err != nil {
				if !k8serrors.IsNotFound(err) {
					fmt.Printf("Error deleting service %s in namespace %s: %v\n", service.Name, service.Namespace, err)
				}
			} else {
				fmt.Printf("Deleted service: %s in namespace: %s\n", service.Name, service.Namespace)
			}
		}
	}

	return nil
}

func (k *K8s) GetCRD(group, version, plural, namespace, name string) (*CRDObject, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	obj, err := k.dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting resource: %v", err)
	}

	// fmt.Printf("Retrieved object: %+v\n", obj) // 添加日志

	status, err := k.GetCRDStatus(obj, group, plural)
	if err != nil {
		return nil, fmt.Errorf("error getting CRD status: %v", err)
	}

	annotations, err := json.Marshal(obj.GetAnnotations())
	if err != nil {
		return nil, fmt.Errorf("error marshalling annotations: %v", err)
	}

	labels, err := json.Marshal(obj.GetLabels())
	if err != nil {
		return nil, fmt.Errorf("error marshalling labels: %v", err)
	}

	spec, err := json.Marshal(obj.Object["spec"])
	if err != nil {
		return nil, fmt.Errorf("error marshalling spec: %v", err)
	}

	statusMore, err := json.Marshal(obj.Object["status"])
	if err != nil {
		return nil, fmt.Errorf("error marshalling status: %v", err)
	}

	createTime := obj.GetCreationTimestamp().Format("2006-01-02 15:04:05")

	crdObj := &CRDObject{
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		Annotations: string(annotations),
		Labels:      string(labels),
		Spec:        string(spec),
		CreateTime:  createTime,
		Status:      status,
		StatusMore:  string(statusMore),
	}

	fmt.Printf("Created CRDObject: %+v\n", crdObj) // 添加日志

	return crdObj, nil
}

func (k *K8s) DeleteCRD(group, version, plural, namespace, name string) error {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	err := k.dynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (k *K8s) CreateCRD(group, version, plural, namespace string, obj map[string]interface{}) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	unstructuredObj := &unstructured.Unstructured{Object: obj}

	createdObj, err := k.dynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create custom resource: %v", err)
	}

	return createdObj, nil
}

func (k *K8s) GetCRDStatus(obj *unstructured.Unstructured, group, plural string) (string, error) {
	switch {
	case plural == "workflows":
		return getWorkflowStatus(obj)
	case plural == "notebooks":
		return getNotebookStatus(obj)
	case plural == "inferenceservices":
		return getInferenceServiceStatus(obj)
	case plural == "jobs" && group == "batch.volcano.sh":
		return getVolcanoJobStatus(obj)
	default:
		return getDefaultStatus(obj)
	}
}

func getDefaultStatus(obj *unstructured.Unstructured) (status string, err error) {
	if obj == nil {
		return "", fmt.Errorf("unstructured object is nil")
	}
	_, exists := obj.Object["status"]
	if !exists {
		// 如果没有 status 字段，返回空字符串而不是错误
		return "", nil
	}

	status, found, err := unstructured.NestedString(obj.Object, "status", "phase")
	if err != nil {
		return "", fmt.Errorf("error getting status.phase: %v", err)
	}
	if found {
		return status, nil
	}

	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return "", fmt.Errorf("error getting status.conditions: %v", err)
	}
	if !found {
		return "", fmt.Errorf("status.conditions not found")
	}

	if len(conditions) == 0 {
		return "", fmt.Errorf("status.conditions is empty")
	}

	lastCondition, ok := conditions[len(conditions)-1].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("last condition is not a map")
	}

	status, found, err = unstructured.NestedString(lastCondition, "type")
	if err != nil {
		return "", fmt.Errorf("error getting condition type: %v", err)
	}
	if !found {
		return "", fmt.Errorf("condition type not found")
	}

	return status, nil
}

func getWorkflowStatus(obj *unstructured.Unstructured) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("unstructured object is nil")
	}

	status, found, err := unstructured.NestedString(obj.Object, "status", "phase")
	if err != nil {
		return "", fmt.Errorf("error getting status.phase: %v", err)
	}
	if !found {
		status = "Unknown"
	}

	nodes, found, err := unstructured.NestedMap(obj.Object, "status", "nodes")
	if err != nil {
		return "", fmt.Errorf("error getting status.nodes: %v", err)
	}
	if !found || len(nodes) == 0 {
		return status, nil
	}

	var lastKey string
	for k := range nodes {
		lastKey = k
	}

	lastNode, ok := nodes[lastKey].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("last node is not a map")
	}

	lastNodePhase, found, err := unstructured.NestedString(lastNode, "phase")
	if err != nil {
		return "", fmt.Errorf("error getting last node phase: %v", err)
	}
	if !found {
		return status, nil
	}

	if lastNodePhase != "Pending" {
		status, _, err = unstructured.NestedString(obj.Object, "status", "phase")
		if err != nil {
			return "", fmt.Errorf("error getting final status.phase: %v", err)
		}
	} else {
		status = lastNodePhase
	}

	return status, nil
}

func getNotebookStatus(obj *unstructured.Unstructured) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("unstructured object is nil")
	}

	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return "", fmt.Errorf("error getting status.conditions: %v", err)
	}
	if !found {
		return "", fmt.Errorf("status.conditions not found")
	}

	if len(conditions) > 0 {
		condition, ok := conditions[0].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("first condition is not a map")
		}
		status, found, err := unstructured.NestedString(condition, "type")
		if err != nil {
			return "", fmt.Errorf("error getting condition type: %v", err)
		}
		if !found {
			return "", fmt.Errorf("condition type not found")
		}
		return status, nil
	}
	return "", fmt.Errorf("no conditions found")
}

func getInferenceServiceStatus(obj *unstructured.Unstructured) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("unstructured object is nil")
	}

	status := "unready"
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return "", fmt.Errorf("error getting status.conditions: %v", err)
	}
	if !found {
		return status, nil
	}

	for _, c := range conditions {
		condition, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		conditionType, found, err := unstructured.NestedString(condition, "type")
		if err != nil {
			return "", fmt.Errorf("error getting condition type: %v", err)
		}
		if !found {
			continue
		}
		conditionStatus, found, err := unstructured.NestedString(condition, "status")
		if err != nil {
			return "", fmt.Errorf("error getting condition status: %v", err)
		}
		if !found {
			continue
		}
		if conditionType == "Ready" && conditionStatus == "True" {
			status = "ready"
			break
		}
	}
	return status, nil
}

func getVolcanoJobStatus(obj *unstructured.Unstructured) (string, error) {
	if obj == nil {
		return "", fmt.Errorf("unstructured object is nil")
	}

	status, found, err := unstructured.NestedString(obj.Object, "status", "state", "phase")
	if err != nil {
		return "", fmt.Errorf("error getting status.state.phase: %v", err)
	}
	if !found {
		status = "unready"
	}
	return status, nil
}

// ServiceExists 检查指定名称的 Service 是否存在
func (k *K8s) ServiceExists(namespace, name string) (bool, error) {
	_, err := k.clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Service 不存在
			return false, nil
		}
		// 发生其他错误
		return false, fmt.Errorf("error checking service existence: %v", err)
	}
	// Service 存在
	return true, nil
}

func (k *K8s) GetPod(namespace, name string) (*corev1.Pod, error) {
	return k.clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdatePod 更新指定的 Pod
func (k *K8s) UpdatePod(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	return k.clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
}

// GetService 获取指定命名空间和名称的 Service
func (k *K8s) GetService(namespace, name string) (*corev1.Service, error) {
	return k.clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// UpdateService 更新指定的 Service
func (k *K8s) UpdateService(namespace string, service *corev1.Service) (*corev1.Service, error) {
	return k.clientset.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
}
