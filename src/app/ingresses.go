package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) controlServiceIngresses(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, ingress, action := v["name"], v["namespace"], v["ingress"], v["action"]
	if err := authObj("service", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "create":
		svc, err := kube.GetService(name, namespace)
		if err != nil {
			return makeStringError(err)
		}
		err = crd.CreateServiceIngress(svc, constant.IdLabel, p.Id,
			a.Args().Namespace(),
			a.Args().Service(),
			a.Args().TLSIssuer(),
			a.Args().BaseDomain())
		if err != nil {
			return makeStringError(err)
		}
	case "delete":
		err = crd.DeleteServiceIngress(a.Args().Namespace(), ingress)
		if err != nil {
			return makeStringError(err)
		}
	}
	return nil
}

func (a *App) watchServiceIngresses(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("service", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.ServiceIngresses(a.Args().Namespace(), name)).ServeHTTP(w, r)
}
