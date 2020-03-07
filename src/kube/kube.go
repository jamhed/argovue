package kube

import (
	"io"
	"os"
	"path/filepath"

	v1client "argovue/client/clientset/versioned"

	fluxv1c "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned"
	fluxv1 "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned/typed/helm.fluxcd.io/v1"

	argovuev1 "argovue/client/clientset/versioned/typed/argovue.io/v1"

	versioned "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return rest.InClusterConfig()
	}
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

func GetClient() (*kubernetes.Clientset, error) {
	config, _ := GetConfig()
	return kubernetes.NewForConfig(config)
}

func GetPodLogs(name, namespace, container string) (io.ReadCloser, error) {
	config, _ := GetConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	podLogOpts := corev1.PodLogOptions{Container: container, Follow: true}
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &podLogOpts)
	return req.Stream()
}

func GetWfClient(wfClientset *versioned.Clientset, namespace string) v1alpha1.WorkflowInterface {
	return wfClientset.ArgoprojV1alpha1().Workflows(namespace)
}

func GetWfClientset() (*versioned.Clientset, error) {
	config, _ := GetConfig()
	return versioned.NewForConfig(config)
}

func GetV1Client(v1Clientset *v1client.Clientset, namespace string) argovuev1.ServiceInterface {
	return v1Clientset.ArgovueV1().Services(namespace)
}

func GetV1Clientset() (*v1client.Clientset, error) {
	config, _ := GetConfig()
	return v1client.NewForConfig(config)
}

func GetFluxV1Client(clientset *fluxv1c.Clientset, namespace string) fluxv1.HelmReleaseInterface {
	return clientset.HelmV1().HelmReleases(namespace)
}

func GetFluxV1Clientset() (*fluxv1c.Clientset, error) {
	config, _ := GetConfig()
	return fluxv1c.NewForConfig(config)
}
