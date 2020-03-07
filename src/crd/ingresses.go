package crd

import (
	"argovue/constant"
	"argovue/kube"
	"argovue/profile"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	extsv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateServiceIngress(svc *corev1.Service, label, owner, ourNamespace, ourServiceName, tlsIssuer, baseDomain string) error {
	host := fmt.Sprintf("%s.%s.%d.svc.cluster.%s", svc.Namespace, svc.Name, svc.Spec.Ports[0].Port, baseDomain)
	ingressName := fmt.Sprintf("%s-%s", svc.Namespace, svc.Name)
	ingress := &extsv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: ourNamespace,
			Annotations: map[string]string{
				"cert-manager.io/cluster-issuer":                 tlsIssuer,
				"nginx.ingress.kubernetes.io/rewrite-target":     "/domain/$1",
				"nginx.ingress.kubernetes.io/proxy-body-size":    "0",
				"nginx.ingress.kubernetes.io/proxy-read-timeout": "600",
				"nginx.ingress.kubernetes.io/proxy-send-timeout": "600",
				constant.OwnerLabel:                              owner,
			},
			Labels: map[string]string{
				constant.ServiceLabel:          svc.Name,
				constant.ServiceNamespaceLabel: svc.Namespace,
				label:                          profile.MaybeHash(label, owner),
			},
		},
		Spec: extsv1.IngressSpec{
			Rules: []extsv1.IngressRule{{
				Host: host,
				IngressRuleValue: extsv1.IngressRuleValue{
					HTTP: &extsv1.HTTPIngressRuleValue{
						Paths: []extsv1.HTTPIngressPath{{
							Backend: extsv1.IngressBackend{ServiceName: ourServiceName, ServicePort: intstr.IntOrString{IntVal: svc.Spec.Ports[0].Port}},
							Path:    "/(.*)",
						}},
					},
				},
			},
			},
			TLS: []extsv1.IngressTLS{
				{
					Hosts:      []string{host},
					SecretName: ingressName,
				},
			},
		},
	}
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	_, err = clientset.ExtensionsV1beta1().Ingresses(ingress.GetNamespace()).Create(ingress)
	return err
}

func DeleteServiceIngress(namespace, name string) error {
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.ExtensionsV1beta1().Ingresses(namespace).Delete(name, opts)
}
