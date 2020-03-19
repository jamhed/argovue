package crd

import (
	"argovue/constant"
	"argovue/kube"
	"argovue/profile"
	"fmt"

	argovuev1 "argovue/apis/argovue.io/v1"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreatePvcDatasource(pvc *corev1.PersistentVolumeClaim, label, owner string) error {
	id := GetIdFromAnnotations("pvc", pvc.Namespace, pvc.Name, constant.DatasourceId)
	datasource := &argovuev1.Datasource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", pvc.Name, id),
			Namespace: pvc.Namespace,
			Annotations: map[string]string{
				constant.OwnerLabel: owner,
			},
			Labels: map[string]string{
				constant.PvcLabel: pvc.Name,
				label:             profile.MaybeHash(label, owner),
			},
		},
		Spec: argovuev1.DatasourceSpec{
			Location: uuid.New().String(),
			Source:   pvc.Name,
		},
	}
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	_, err = clientset.ArgovueV1().Datasources(pvc.Namespace).Create(datasource)
	return err
}

func DeletePvcDatasource(namespace, name string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.ArgovueV1().Datasources(namespace).Delete(name, opts)
}

func CreateDatasourcePvc(datasource *argovuev1.Datasource, label, owner string) error {
	id := GetIdFromAnnotations("datasource", datasource.Namespace, datasource.Name, constant.PvcId)
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", datasource.Name, id),
			Namespace: datasource.Namespace,
			Annotations: map[string]string{
				constant.OwnerLabel: owner,
			},
			Labels: map[string]string{
				constant.DatasourceLabel: datasource.Name,
				label:                    profile.MaybeHash(label, owner),
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources:   corev1.ResourceRequirements{Requests: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("1Gi")}},
		},
	}
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	_, err = clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(pvc)
	return err
}

func DeleteDatasourcePvc(namespace, name string) error {
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(name, opts)
}
