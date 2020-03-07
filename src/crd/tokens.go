package crd

import (
	"argovue/constant"
	"argovue/kube"
	"argovue/profile"
	"fmt"

	argovuev1 "argovue/apis/argovue.io/v1"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateServiceToken(svc *corev1.Service, label, owner string) error {
	id := GetIdFromAnnotations("service", svc.Namespace, svc.Name, constant.TokenId)
	token := &argovuev1.Token{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", svc.Name, id),
			Namespace: svc.Namespace,
			Annotations: map[string]string{
				constant.OwnerLabel: owner,
			},
			Labels: map[string]string{
				constant.ServiceLabel: svc.Name,
				label:                 profile.MaybeHash(label, owner),
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "v1", Kind: "Service", Name: svc.Name, UID: svc.UID}},
		},
		Spec: argovuev1.TokenSpec{
			Value:       uuid.New().String(),
			Description: "",
		},
	}
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	_, err = clientset.ArgovueV1().Tokens(svc.Namespace).Create(token)
	return err
}

func DeleteServiceToken(namespace, name string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.ArgovueV1().Tokens(namespace).Delete(name, opts)
}
