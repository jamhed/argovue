package kube

import (
	"fmt"

	"argovue/constant"

	argovuev1 "argovue/apis/argovue.io/v1"
	v1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	fluxv1 "github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetByKind(kind, name, namespace string) (metav1.Object, error) {
	switch kind {
	case "service":
		return GetService(name, namespace)
	case "pod":
		return GetPod(name, namespace)
	case "pvc":
		return GetPvc(name, namespace)
	case "deployment":
		return GetDeployment(name, namespace)
	case "workflow":
		return GetWorkflow(name, namespace)
	case "ingress":
		return GetIngress(name, namespace)
	case "dataset":
		return GetDataset(name, namespace)
	case "job":
		return GetJob(name, namespace)
	case "argovue":
		return GetArgovueService(name, namespace)
	case "helmrelease":
		return GetHelmRelease(name, namespace)
	default:
		return nil, fmt.Errorf("Unknown kubernetes kind %s", kind)
	}
}

func GetArgovueService(name, namespace string) (*argovuev1.Service, error) {
	clientset, err := GetV1Clientset()
	if err != nil {
		return nil, err
	}
	return clientset.ArgovueV1().Services(namespace).Get(name, metav1.GetOptions{})
}

func GetDataset(name, namespace string) (*argovuev1.Dataset, error) {
	clientset, err := GetV1Clientset()
	if err != nil {
		return nil, err
	}
	return clientset.ArgovueV1().Datasets(namespace).Get(name, metav1.GetOptions{})
}

func GetJob(name, namespace string) (*batchv1.Job, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})
}

func GetWorkflow(name, namespace string) (*v1alpha1.Workflow, error) {
	clientset, err := GetWfClientset()
	if err != nil {
		return nil, err
	}
	return clientset.ArgoprojV1alpha1().Workflows(namespace).Get(name, metav1.GetOptions{})
}

func GetService(name, namespace string) (*corev1.Service, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
}

func GetServiceTokens(name, namespace string) (*argovuev1.TokenList, error) {
	clientset, err := GetV1Clientset()
	if err != nil {
		return nil, err
	}
	return clientset.ArgovueV1().Tokens(namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", constant.ServiceLabel, name),
	})
}

func GetPod(name, namespace string) (*corev1.Pod, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

func GetPvc(name, namespace string) (*corev1.PersistentVolumeClaim, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})
}

func GetIngress(name, namespace string) (*extsv1beta1.Ingress, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.ExtensionsV1beta1().Ingresses(namespace).Get(name, metav1.GetOptions{})
}

func GetDeployment(name, namespace string) (*appsv1.Deployment, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}

func GetHelmRelease(name, namespace string) (*fluxv1.HelmRelease, error) {
	clientset, err := GetFluxV1Clientset()
	if err != nil {
		return nil, err
	}
	return clientset.HelmV1().HelmReleases(namespace).Get(name, metav1.GetOptions{})
}
