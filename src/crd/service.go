package crd

import (
	"fmt"
	"strconv"

	"argovue/constant"
	"argovue/kube"
	"argovue/profile"

	argovuev1 "argovue/apis/argovue.io/v1"
	fluxv1 "github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"

	wfv1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func int32Ptr(i int32) *int32 { return &i }

func ensureArgovueValues(release *fluxv1.HelmRelease) *fluxv1.HelmRelease {
	if release.Spec.Values == nil {
		release.Spec.Values = make(map[string]interface{})
	}
	av, ok := release.Spec.Values["argovue"].(map[string]interface{})
	if !ok {
		av = make(map[string]interface{})
	}
	labels, ok := (av["labels"]).(map[string]string)
	if !ok {
		labels = make(map[string]string)
	}
	av["labels"] = labels
	release.Spec.Values["argovue"] = av
	return release
}

func addArgovueLabel(release *fluxv1.HelmRelease, label, value string) {
	release.Spec.Values["argovue"].(map[string]interface{})["labels"].(map[string]string)[label] = value
}

func addArgovueValue(release *fluxv1.HelmRelease, key string, value interface{}) {
	release.Spec.Values["argovue"].(map[string]interface{})[key] = value
}

func makeRelease(s *argovuev1.Service, namespace, label, owner, releaseName string) *fluxv1.HelmRelease {
	release := ensureArgovueValues(&fluxv1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      releaseName,
			Namespace: namespace,
			Annotations: map[string]string{
				constant.OwnerLabel: owner,
			},
			Labels: map[string]string{
				constant.ServiceLabel: s.Name,
				label:                 profile.MaybeHash(label, owner),
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Service", Name: s.Name, UID: s.UID}},
		},
		Spec: s.Spec.HelmRelease,
	})
	release.Spec.ReleaseName = releaseName
	addArgovueLabel(release, label, profile.MaybeHash(label, owner))
	addArgovueValue(release, "owner", owner)
	addArgovueValue(release, "ownerLabel", constant.OwnerLabel)
	addArgovueValue(release, "baseurl", fmt.Sprintf("/proxy/%s/%s/%d", namespace, releaseName, 80))
	return release
}

func deployRelease(release *fluxv1.HelmRelease) error {
	clientset, err := kube.GetFluxV1Clientset()
	if err != nil {
		return err
	}
	_, err = clientset.HelmV1().HelmReleases(release.GetNamespace()).Create(release)
	return err
}

func Deploy(s *argovuev1.Service, label, owner string, canonical string, input []argovuev1.InputValue) error {
	var releaseName string
	if canonical == "true" {
		releaseName = s.Name
	} else {
		releaseName = fmt.Sprintf("%s-%s", s.Name, GetIdFromAnnotations("argovue", s.Namespace, s.Name, constant.InstanceId))
	}
	release := makeRelease(s, s.Namespace, label, owner, releaseName)
	addArgovueValue(release, "input", input)
	return deployRelease(release)
}

func DeployFilebrowser(wf *wfv1alpha1.Workflow, namespace, argovueReleaseName, owner string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	filebrowser, err := clientset.ArgovueV1().Services(namespace).Get(fmt.Sprintf("%s-filebrowser", argovueReleaseName), metav1.GetOptions{})
	if err != nil {
		return err
	}
	releaseName := fmt.Sprintf("%s-%s", wf.Name, GetIdFromAnnotations("workflow", wf.Namespace, wf.Name, constant.InstanceId))
	release := makeRelease(filebrowser, wf.Namespace, constant.IdLabel, owner, releaseName)
	volumes := []map[string]string{}
	for _, pvc := range wf.Status.PersistentVolumeClaims {
		volumes = append(volumes, map[string]string{"name": pvc.Name, "claim": pvc.PersistentVolumeClaim.ClaimName})
	}
	release.ObjectMeta.OwnerReferences =
		[]metav1.OwnerReference{{APIVersion: "argoproj.io/v1alpha1", Kind: "Workflow", Name: wf.Name, UID: wf.UID}}
	release.ObjectMeta.Labels[constant.WorkflowLabel] = wf.Name
	addArgovueLabel(release, constant.WorkflowLabel, wf.Name)
	release.Spec.Values["volumes"] = volumes
	return deployRelease(release)
}

func DeleteInstance(namespace, name string) error {
	clientset, err := kube.GetFluxV1Clientset()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.HelmV1().HelmReleases(namespace).Delete(name, opts)
}

func GetWorkflowFilebrowserNames(wf *wfv1alpha1.Workflow) (re []string) {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		log.Errorf("Can't get clientset, error:%s", err)
		return
	}
	iface := clientset.ArgovueV1().Services(wf.Namespace)

	list, err := iface.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,app.kubernetes.io/name=%s", constant.WorkflowLabel, wf.Name, "filebrowser")})
	if err != nil {
		return
	}
	for _, svc := range list.Items {
		re = append(re, svc.GetName())
	}
	return
}

func GetIdFromAnnotations(kind, namespace, name, label string) string {
	client := kube.ByKind(kind, namespace)
	obj, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get object %s/%s/%s, error:%s", kind, namespace, name, err)
		return "0"
	}
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	id, ok := annotations[label]
	if !ok {
		id = "1"
	} else {
		idI, err := strconv.Atoi(id)
		if err != nil {
			idI = 1
		}
		id = strconv.Itoa(idI + 1)
	}
	annotations[label] = id
	obj.SetAnnotations(annotations)
	_, err = client.Update(obj, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("Can't update annotation %s/%s/%s, error:%s", kind, namespace, name, err)
	}
	return id
}
