package services

import (
	"path/filepath"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreatePod(t *testing.T) {
	// Create a mock Kubernetes clientset
	kubeconfig := filepath.Join(".", "config")

	// Build Kubernetes configuration
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating Kubernetes client: %v", err)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating dynamic client: %v", err)
	}

	k8s := &K8s{
		clientset:     clientset,
		dynamicClient: dynamicClient,
	}

	// Set up Pod parameters
	namespace := "jupyter"
	name := "test-pod"
	labels := map[string]string{
		"app":      "test-pod",
		"pod-type": "notebook",
		"user":     "testuser",
	}
	command := []string{"sh", "-c"}
	args := []string{`jupyter lab --notebook-dir=/mnt/testuser --ip=0.0.0.0 --no-browser --allow-root --port=3000 --NotebookApp.token='' --NotebookApp.password='' --ServerApp.disable_check_xsrf=True --NotebookApp.allow_origin='*' --NotebookApp.base_url=/notebook/jupyter/testuser/ --NotebookApp.tornado_settings='{"headers": {"Content-Security-Policy": "frame-ancestors * 'self' "}}'`}

	// Create Pod object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    name,
					Image:   "jupyter/base-notebook:latest",
					Command: command,
					Args:    args,
					VolumeMounts: []corev1.VolumeMount{
						{Name: "kubeflow-user-workspace", MountPath: "/mnt/testuser", SubPath: "testuser"},
						{Name: "kubeflow-archives", MountPath: "/archives/testuser", SubPath: "testuser"},
						{Name: "tz-config", MountPath: "/etc/localtime"},
						{Name: "dshm", MountPath: "/dev/shm"},
					},
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("4Gi"),
							corev1.ResourceCPU:    resource.MustParse("2"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("4Gi"),
							corev1.ResourceCPU:    resource.MustParse("2"),
						},
					},
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "kubeflow-user-workspace",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "kubeflow-user-workspace",
						},
					},
				},
				{
					Name: "kubeflow-archives",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: "kubeflow-archives",
						},
					},
				},
				{
					Name: "tz-config",
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/usr/share/zoneinfo/Asia/Shanghai",
						},
					},
				},
				{
					Name: "dshm",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMediumMemory,
						},
					},
				},
			},
			RestartPolicy:      corev1.RestartPolicyNever,
			NodeSelector:       map[string]string{"notebook": "true"},
			ImagePullSecrets:   []corev1.LocalObjectReference{{Name: "hubsecret"}},
			ServiceAccountName: "default",
			SchedulerName:      "default-scheduler",
		},
	}

	// Call CreatePod function
	createdPod, err := k8s.CreatePod(namespace, pod)
	if err != nil {
		t.Fatalf("Failed to create pod: %v", err)
	}

	// Verify created Pod
	if createdPod.Name != name {
		t.Errorf("Expected pod name %s, got %s", name, createdPod.Name)
	}
	if createdPod.Namespace != namespace {
		t.Errorf("Expected namespace %s, got %s", namespace, createdPod.Namespace)
	}
	if !reflect.DeepEqual(createdPod.Labels, labels) {
		t.Errorf("Labels do not match. Expected %v, got %v", labels, createdPod.Labels)
	}
	if createdPod.Spec.Containers[0].Image != "jupyter/base-notebook:latest" {
		t.Errorf("Expected image %s, got %s", "jupyter/base-notebook:latest", createdPod.Spec.Containers[0].Image)
	}
	if !reflect.DeepEqual(createdPod.Spec.Containers[0].Command, command) {
		t.Errorf("Command does not match. Expected %v, got %v", command, createdPod.Spec.Containers[0].Command)
	}
	if !reflect.DeepEqual(createdPod.Spec.Containers[0].Args, args) {
		t.Errorf("Args do not match. Expected %v, got %v", args, createdPod.Spec.Containers[0].Args)
	}
	if createdPod.Spec.RestartPolicy != corev1.RestartPolicyNever {
		t.Errorf("Expected restart policy %v, got %v", corev1.RestartPolicyNever, createdPod.Spec.RestartPolicy)
	}
	if !reflect.DeepEqual(createdPod.Spec.NodeSelector, map[string]string{"notebook": "true"}) {
		t.Errorf("Node selector does not match. Expected %v, got %v", map[string]string{"notebook": "true"}, createdPod.Spec.NodeSelector)
	}
	if createdPod.Spec.Containers[0].ImagePullPolicy != corev1.PullIfNotPresent {
		t.Errorf("Image pull policy does not match. Expected %v, got %v", corev1.PullIfNotPresent, createdPod.Spec.Containers[0].ImagePullPolicy)
	}
	if !reflect.DeepEqual(createdPod.Spec.ImagePullSecrets, []corev1.LocalObjectReference{{Name: "hubsecret"}}) {
		t.Errorf("Image pull secrets do not match. Expected %v, got %v", []corev1.LocalObjectReference{{Name: "hubsecret"}}, createdPod.Spec.ImagePullSecrets)
	}
	if createdPod.Spec.SchedulerName != "default-scheduler" {
		t.Errorf("Scheduler name does not match. Expected %s, got %s", "default-scheduler", createdPod.Spec.SchedulerName)
	}
}
