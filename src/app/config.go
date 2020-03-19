package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/kube"
	"argovue/msg"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"
)

func (a *App) updateGroups(obj interface{}) {
	cfg, err := crd.TypecastConfig(obj)
	if err != nil {
		return
	}
	for _, i := range cfg.Spec.Groups {
		a.groups[i.Oidc] = i.Kubernetes
	}
}

func (a *App) ListenForServices() {
	for {
		select {
		case m, ok := <-a.services.Notifier():
			if !ok {
				log.Debugf("Services stop")
				return
			}
			if m.Action == "delete" {
				svc := m.Content.(v1.Object)
				name := svc.GetName()
				namespace := svc.GetNamespace()
				labelSelector := fmt.Sprintf("%s=%s,%s=%s", constant.ServiceLabel, name, constant.ServiceNamespaceLabel, namespace)
				list, err := kube.GetIngressesByLabel(a.Args().Namespace(), labelSelector)
				if err != nil {
					log.Errorf("Error getting ingress objects for service %s/%s %s", namespace, name, err)
					return
				}
				for _, ingress := range list.Items {
					err := kube.DeleteIngress(ingress.GetName(), a.Args().Namespace())
					if err != nil {
						log.Errorf("Error deleting ingress:%s object for service %s/%s %s", ingress.GetName(), namespace, name, err)
					} else {
						log.Infof("Clean up ingress:%s", ingress.GetName())
					}
				}
			}
		}
	}
}

func (a *App) ListenForConfig() {
	for {
		select {
		case msg, ok := <-a.config.Notifier():
			if !ok {
				log.Debugf("Config stop")
				return
			}
			a.updateConfig(msg)
		}
	}
}

func (a *App) updateConfig(msg *msg.Msg) {
	a.groups = make(map[string]string)
	switch msg.Action {
	case "update":
		a.updateGroups(msg.Content)
	case "add":
		a.updateGroups(msg.Content)
	default:
	}
	log.Debugf("App: configured groups %s", a.groups)
}
