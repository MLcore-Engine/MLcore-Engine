package services

import (
	"context"
	"fmt"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// PyTorchJobConfig defines the configuration for a PyTorch training job
type PyTorchJobConfig struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Image           string            `json:"image"`
	ImagePullPolicy string            `json:"image_pull_policy"`
	RestartPolicy   string            `json:"restart_policy"`
	Command         []string          `json:"command"`
	Args            []string          `json:"args,omitempty"`
	MasterReplicas  int32             `json:"master_replicas"`
	WorkerReplicas  int32             `json:"worker_replicas"`
	GPUsPerNode     int64             `json:"gpus_per_node"`
	CPULimit        string            `json:"cpu_limit"`
	MemoryLimit     string            `json:"memory_limit"`
	NodeSelector    map[string]string `json:"node_selector"`
	Env             []EnvVar          `json:"env"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (k *K8s) CreatePyTorchJob(namespace string, config PyTorchJobConfig) (*unstructured.Unstructured, error) {

	ctx := context.Background()
	if config.Name == "" {
		return nil, fmt.Errorf("job name cannot be empty")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be empty")
	}

	pytorchJob := map[string]interface{}{
		"apiVersion": "kubeflow.org/v1",
		"kind":       "PyTorchJob",
		"metadata": map[string]interface{}{
			"name":      config.Name,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"pytorchReplicaSpecs": map[string]interface{}{
				"Master": createReplicaSpec(config, "Master"),
				"Worker": createReplicaSpec(config, "Worker"),
			},
		},
	}

	gvr := schema.GroupVersionResource{
		Group:    "kubeflow.org",
		Version:  "v1",
		Resource: "pytorchjobs",
	}

	// try delete existing job
	err := k.dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, config.Name, metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		fmt.Printf("Error deleting existing PyTorch job: %v\n", err)
	}

	createdJob, err := k.dynamicClient.Resource(gvr).Namespace(namespace).Create(
		ctx,
		&unstructured.Unstructured{Object: pytorchJob},
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PyTorch job: %v", err)
	}

	time.Sleep(time.Second)

	return createdJob, nil
}

func (k *K8s) DeletePyTorchJob(namespace string, name string) error {
	ctx := context.Background()

	gvr := schema.GroupVersionResource{
		Group:    "kubeflow.org",
		Version:  "v1",
		Resource: "pytorchjobs",
	}

	err := k.dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete PyTorch job: %v", err)
	}

	return nil
}

func (k *K8s) GetPyTorchJob(namespace string, name string) (*unstructured.Unstructured, error) {
	ctx := context.Background()

	gvr := schema.GroupVersionResource{
		Group:    "kubeflow.org",
		Version:  "v1",
		Resource: "pytorchjobs",
	}

	job, err := k.dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PyTorch job: %v", err)
	}

	return job, nil
}

func createReplicaSpec(config PyTorchJobConfig, replicaType string) map[string]interface{} {
	var replicas int32
	if replicaType == "Master" {
		replicas = config.MasterReplicas
	} else if replicaType == "Worker" {
		replicas = config.WorkerReplicas
	}

	return map[string]interface{}{
		"replicas":      replicas,
		"restartPolicy": config.RestartPolicy,
		"template": map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []map[string]interface{}{
					{
						"name":            "pytorch",
						"image":           config.Image,
						"imagePullPolicy": config.ImagePullPolicy,
						"command":         config.Command,
						"args":            config.Args,
						"resources": map[string]interface{}{
							"limits": map[string]interface{}{
								"cpu":            config.CPULimit,
								"memory":         config.MemoryLimit,
								"nvidia.com/gpu": config.GPUsPerNode,
							},
							"requests": map[string]interface{}{
								"cpu":            config.CPULimit,
								"memory":         config.MemoryLimit,
								"nvidia.com/gpu": config.GPUsPerNode,
							},
						},
						// "env": createEnvVars(config.Env, replicaType),
					},
				},
				"nodeSelector": config.NodeSelector,
			},
		},
	}
}

func createEnvVars(envVars []EnvVar, replicaType string) []map[string]interface{} {

	env := []map[string]interface{}{}

	// add custom env vars
	for _, e := range envVars {
		env = append(env, map[string]interface{}{
			"name":  e.Name,
			"value": e.Value,
		})
	}

	// add distributed training env vars
	env = append(env, map[string]interface{}{
		"name": "MASTER_ADDR",
		"valueFrom": map[string]interface{}{
			"fieldRef": map[string]interface{}{
				"fieldPath": "status.podIP",
			},
		},
	})
	if replicaType == "Master" {
		env = append(env, map[string]interface{}{
			"name":  "MASTER_PORT",
			"value": "29500",
		})
	}

	return env
}
