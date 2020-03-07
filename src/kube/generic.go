package kube

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Dynamic struct {
	group    string
	version  string
	resource string
	gvr      schema.GroupVersionResource
}

func ArgovueService() *Dynamic {
	return new("argovue.io", "v1", "services")
}

func ArgovueDataset() *Dynamic {
	return new("argovue.io", "v1", "datasets")
}

func ArgoWorkflow() *Dynamic {
	return new("argoproj.io", "v1alpha1", "workflows")
}

func Service() *Dynamic {
	return new("", "v1", "services")
}

func Pvc() *Dynamic {
	return new("", "v1", "persistentvolumeclaims")
}

func ByKind(name, namespace string) dynamic.ResourceInterface {
	switch name {
	case "workflow":
		return ArgoWorkflow().Namespace(namespace)
	case "argovue":
		return ArgovueService().Namespace(namespace)
	case "service":
		return Service().Namespace(namespace)
	case "pvc":
		return Pvc().Namespace(namespace)
	case "dataset":
		return ArgovueDataset().Namespace(namespace)
	}
	return nil
}

func new(group, version, resource string) *Dynamic {
	return &Dynamic{
		group:    group,
		version:  version,
		resource: resource,
		gvr:      schema.GroupVersionResource{Group: group, Version: version, Resource: resource},
	}
}

func getDynamicClient() dynamic.Interface {
	config, _ := GetConfig()
	client, _ := dynamic.NewForConfig(config)
	return client
}

func (g *Dynamic) Namespace(namespace string) dynamic.ResourceInterface {
	client := getDynamicClient()
	return client.Resource(g.gvr).Namespace(namespace)
}
