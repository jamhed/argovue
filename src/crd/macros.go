package crd

import (
	"argovue/constant"
	"encoding/json"
	"fmt"

	v1 "argovue/apis/argovue.io/v1"
)

func GetByKind(kind, namespace, name string) (crd *Crd, err error) {
	err = nil
	switch kind {
	case "pod":
		crd = New("", "v1", "pods")
	case "job":
		crd = New("batch", "v1", "jobs")
	case "pvc":
		crd = New("", "v1", "persistentvolumeclaims")
	case "service":
		crd = New("", "v1", "services")
	case "workflow":
		crd = New("argoproj.io", "v1alpha1", "workflows")
	case "argovue":
		crd = New("argovue.io", "v1", "services")
	case "token":
		crd = New("argovue.io", "v1", "tokens")
	case "ingress":
		crd = New("extensions", "v1beta1", "ingresses")
	case "datasource":
		crd = New("argovue.io", "v1", "datasources")
	case "helmrelease":
		crd = New("helm.fluxcd.io", "v1", "helmreleases")
	default:
		return nil, fmt.Errorf("Can't create crd by kind:%s", kind)
	}
	crd.SetFieldSelector(fmt.Sprintf("metadata.name=%s,metadata.namespace=%s", name, namespace))
	return
}

func WorkflowPods(namespace, name, pod string) *Crd {
	return New("", "v1", "pods").
		SetLabelSelector(fmt.Sprintf("%s=%s", constant.WorkflowLabel, name)).
		SetFieldSelector(fmt.Sprintf("metadata.namespace=%s,metadata.name=%s", namespace, pod))
}

func WorkflowServices(namespace, name string) *Crd {
	return New("helm.fluxcd.io", "v1", "helmreleases").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector(fmt.Sprintf("%s=%s", constant.WorkflowLabel, name))
}

func WorkflowMounts(namespace, name string) *Crd {
	return New("", "v1", "services").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector(fmt.Sprintf("%s=%s", constant.WorkflowLabel, name))
}

func Workflow(namespace, name string) *Crd {
	return New("argoproj.io", "v1alpha1", "workflows").
		SetFieldSelector(fmt.Sprintf("metadata.namespace=%s,metadata.name=%s", namespace, name))
}

func Catalogue(namespace, name string) *Crd {
	return New("argovue.io", "v1", "services").
		SetFieldSelector(fmt.Sprintf("metadata.namespace=%s,metadata.name=%s", namespace, name))
}

func CatalogueInstances(namespace, name string) *Crd {
	return New("helm.fluxcd.io", "v1", "helmreleases").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("service.argovue.io/name=" + name)
}

func ServiceTokens(namespace, name string) *Crd {
	return New("argovue.io", "v1", "tokens").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("service.argovue.io/name=" + name)
}

func PvcDatasources(namespace, name string) *Crd {
	return New("argovue.io", "v1", "datasources").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("pvc.argovue.io/name=" + name)
}

func DatasourcePvcs(namespace, name string) *Crd {
	return New("", "v1", "persistentvolumeclaims").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector(constant.DatasourceLabel + "=" + name)
}

func PvcMounts(namespace, name string) *Crd {
	return New("", "v1", "services").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("pvc.argovue.io/name=" + name)
}

func DatasourceUploads(namespace, name string) *Crd {
	return New("batch", "v1", "jobs").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("datasource.argovue.io/name=" + name)
}

func JobPods(namespace, name string) *Crd {
	return New("", "v1", "pods").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("job-name=" + name)
}

func ServiceIngresses(namespace, name string) *Crd {
	return New("extensions", "v1beta1", "ingresses").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("service.argovue.io/name=" + name)
}

func CatalogueResources(namespace, name string) *Crd {
	return New("", "v1", "pods").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("app.kubernetes.io/name=" + name)
}

func CatalogueWorkflows(namespace, name string) *Crd {
	return New("argoproj.io", "v1alpha1", "workflows").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstancePods(namespace, name string) *Crd {
	return New("", "v1", "pods").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstanceServices(namespace, name string) *Crd {
	return New("", "v1", "services").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstancePvcs(namespace, name string) *Crd {
	return New("", "v1", "persistentvolumeclaims").
		SetFieldSelector("metadata.namespace=" + namespace).
		SetLabelSelector("app.kubernetes.io/instance=" + name)
}

func CatalogueInstanceIngresses(namespace, name string) *Crd {
	return New("extensions", "v1beta1", "ingresses").
		SetLabelSelector(fmt.Sprintf("%s=%s,%s=%s", constant.ServiceLabel, name, constant.ServiceNamespaceLabel, namespace))
}

func CatalogueInstance(namespace, name, instance string) *Crd {
	return New("helm.fluxcd.io", "v1", "helmreleases").
		SetLabelSelector("service.argovue.io/name=" + name).
		SetFieldSelector(fmt.Sprintf("metadata.namespace=%s,metadata.name=%s", namespace, instance))
}

func Typecast(thing interface{}) (*v1.Service, error) {
	if thing == nil {
		return nil, fmt.Errorf("Service typecast nil input")
	}
	buf, err := json.Marshal(thing)
	if err != nil {
		return nil, err
	}
	svc := new(v1.Service)
	err = json.Unmarshal(buf, svc)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func TypecastConfig(thing interface{}) (*v1.AppConfig, error) {
	if thing == nil {
		return nil, fmt.Errorf("Service typecast nil input")
	}
	buf, err := json.Marshal(thing)
	if err != nil {
		return nil, err
	}
	cfg := new(v1.AppConfig)
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
