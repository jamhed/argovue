package crd

import (
	"argovue/args"
	"argovue/constant"
	"argovue/kube"
	"argovue/profile"
	"fmt"

	argovuev1 "argovue/apis/argovue.io/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SyncDatasetPvc(dataset *argovuev1.Dataset, pvc *v1.PersistentVolumeClaim, label, owner string, params args.RcloneParams) error {
	id := GetIdFromAnnotations("dataset", dataset.Namespace, dataset.Name, constant.JobId)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("sync-ds-pvc-%s-%s", dataset.Name, id),
			Namespace:   dataset.Namespace,
			Annotations: map[string]string{constant.OwnerLabel: owner},
			Labels: map[string]string{
				constant.DatasetLabel: dataset.Name,
				label:                 profile.MaybeHash(label, owner),
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Dataset", Name: dataset.Name, UID: dataset.UID}},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("%s-%s", dataset.Name, id),
					Namespace:   dataset.Namespace,
					Annotations: map[string]string{constant.OwnerLabel: owner},
					Labels: map[string]string{
						constant.DatasetLabel: dataset.Name,
						label:                 profile.MaybeHash(label, owner),
					},
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{{
						Name:    "rclone",
						Image:   params.Image,
						Command: []string{"rclone", "-v", "sync", "S3:" + params.Bucket + "/" + dataset.Spec.Location, "/pvc"},
						Env: []v1.EnvVar{
							{Name: "RCLONE_CONFIG_S3_ENDPOINT", Value: params.Endpoint},
							{Name: "RCLONE_CONFIG_S3_REGION", Value: params.Region},
							{Name: "RCLONE_CONFIG_S3_ACCESS_KEY_ID", Value: params.Key},
							{Name: "RCLONE_CONFIG_S3_SECRET_ACCESS_KEY", Value: params.Secret},
							{Name: "RCLONE_S3_SESSION_TOKEN", Value: params.Session},
						},
						VolumeMounts: []v1.VolumeMount{{Name: "pvc", MountPath: "/pvc", ReadOnly: false}},
					},
					},
					Volumes: []v1.Volume{{
						Name: "pvc",
						VolumeSource: v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvc.Name,
							ReadOnly:  false,
						}},
					}},
				},
			},
		},
	}
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	_, err = clientset.BatchV1().Jobs(dataset.Namespace).Create(job)
	return err
}

func SyncPvcDataset(dataset *argovuev1.Dataset, label, owner string, params args.RcloneParams) error {
	id := GetIdFromAnnotations("dataset", dataset.Namespace, dataset.Name, constant.JobId)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("sync-pvc-ds-%s-%s", dataset.Name, id),
			Namespace:   dataset.Namespace,
			Annotations: map[string]string{constant.OwnerLabel: owner},
			Labels: map[string]string{
				constant.DatasetLabel: dataset.Name,
				label:                 profile.MaybeHash(label, owner),
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Dataset", Name: dataset.Name, UID: dataset.UID}},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("%s-%s", dataset.Name, id),
					Namespace:   dataset.Namespace,
					Annotations: map[string]string{constant.OwnerLabel: owner},
					Labels: map[string]string{
						constant.DatasetLabel: dataset.Name,
						label:                 profile.MaybeHash(label, owner),
					},
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{{
						Name:    "rclone",
						Image:   params.Image,
						Command: []string{"rclone", "-v", "sync", "/pvc", "S3:" + params.Bucket + "/" + dataset.Spec.Location},
						Env: []v1.EnvVar{
							{Name: "RCLONE_CONFIG_S3_ENDPOINT", Value: params.Endpoint},
							{Name: "RCLONE_CONFIG_S3_REGION", Value: params.Region},
							{Name: "RCLONE_CONFIG_S3_ACCESS_KEY_ID", Value: params.Key},
							{Name: "RCLONE_CONFIG_S3_SECRET_ACCESS_KEY", Value: params.Secret},
							{Name: "RCLONE_S3_SESSION_TOKEN", Value: params.Session},
						},
						VolumeMounts: []v1.VolumeMount{{Name: "pvc", MountPath: "/pvc", ReadOnly: true}},
					},
					},
					Volumes: []v1.Volume{{
						Name: "pvc",
						VolumeSource: v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: dataset.Spec.Source,
							ReadOnly:  true,
						}},
					}},
				},
			},
		},
	}
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	_, err = clientset.BatchV1().Jobs(dataset.Namespace).Create(job)
	return err
}

func DeleteDatasetSync(namespace, name string) error {
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.BatchV1().Jobs(namespace).Delete(name, opts)
}
