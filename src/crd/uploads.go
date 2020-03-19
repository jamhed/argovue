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

func SyncDatasourcePvc(datasource *argovuev1.Datasource, pvc *v1.PersistentVolumeClaim, label, owner string, params args.RcloneParams) error {
	id := GetIdFromAnnotations("datasource", datasource.Namespace, datasource.Name, constant.JobId)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("sync-ds-pvc-%s-%s", datasource.Name, id),
			Namespace:   datasource.Namespace,
			Annotations: map[string]string{constant.OwnerLabel: owner},
			Labels: map[string]string{
				constant.DatasourceLabel: datasource.Name,
				label:                    profile.MaybeHash(label, owner),
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Datasource", Name: datasource.Name, UID: datasource.UID}},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("%s-%s", datasource.Name, id),
					Namespace:   datasource.Namespace,
					Annotations: map[string]string{constant.OwnerLabel: owner},
					Labels: map[string]string{
						constant.DatasourceLabel: datasource.Name,
						label:                    profile.MaybeHash(label, owner),
					},
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{{
						Name:    "rclone",
						Image:   params.Image,
						Command: []string{"rclone", "-v", "sync", "S3:" + params.Bucket + "/" + datasource.Spec.Location, "/pvc"},
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
	_, err = clientset.BatchV1().Jobs(datasource.Namespace).Create(job)
	return err
}

func SyncPvcDatasource(datasource *argovuev1.Datasource, label, owner string, params args.RcloneParams) error {
	id := GetIdFromAnnotations("datasource", datasource.Namespace, datasource.Name, constant.JobId)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("sync-pvc-ds-%s-%s", datasource.Name, id),
			Namespace:   datasource.Namespace,
			Annotations: map[string]string{constant.OwnerLabel: owner},
			Labels: map[string]string{
				constant.DatasourceLabel: datasource.Name,
				label:                    profile.MaybeHash(label, owner),
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Datasource", Name: datasource.Name, UID: datasource.UID}},
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("%s-%s", datasource.Name, id),
					Namespace:   datasource.Namespace,
					Annotations: map[string]string{constant.OwnerLabel: owner},
					Labels: map[string]string{
						constant.DatasourceLabel: datasource.Name,
						label:                    profile.MaybeHash(label, owner),
					},
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{{
						Name:    "rclone",
						Image:   params.Image,
						Command: []string{"rclone", "-v", "sync", "/pvc", "S3:" + params.Bucket + "/" + datasource.Spec.Location},
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
							ClaimName: datasource.Spec.Source,
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
	_, err = clientset.BatchV1().Jobs(datasource.Namespace).Create(job)
	return err
}

func DeleteDatasourceSync(namespace, name string) error {
	clientset, err := kube.GetClient()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.BatchV1().Jobs(namespace).Delete(name, opts)
}
