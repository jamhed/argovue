package crd

import (
	"argovue/constant"
	"argovue/kube"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePvcMount(pvc *corev1.PersistentVolumeClaim, namespace, argovueReleaseName, label, owner, canonical string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	filebrowser, err := clientset.ArgovueV1().Services(namespace).Get(fmt.Sprintf("%s-filebrowser", argovueReleaseName), metav1.GetOptions{})
	if err != nil {
		return err
	}
	var releaseName string
	if canonical == "true" {
		releaseName = fmt.Sprintf("mount-%s", pvc.Name)
	} else {
		releaseName = fmt.Sprintf("mount-%s-%s", pvc.Name, GetIdFromAnnotations("pvc", pvc.Namespace, pvc.Name, constant.InstanceId))
	}
	release := makeRelease(filebrowser, pvc.Namespace, label, owner, releaseName)
	volumes := []map[string]string{{"name": pvc.Name, "claim": pvc.Name}}
	release.ObjectMeta.OwnerReferences =
		[]metav1.OwnerReference{{APIVersion: "v1", Kind: "PersistentVolumeClaim", Name: pvc.Name, UID: pvc.UID}}
	release.ObjectMeta.Labels[constant.PvcLabel] = pvc.Name
	addArgovueLabel(release, constant.PvcLabel, pvc.Name)
	release.Spec.Values["volumes"] = volumes
	return deployRelease(release)
}

func DeletePvcMount(namespace, name string) error {
	return nil
}
