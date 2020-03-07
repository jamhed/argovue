package app

import (
	"argovue/crd"
	"argovue/msg"

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
