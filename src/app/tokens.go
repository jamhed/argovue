package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) controlServiceTokens(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, token, action := v["name"], v["namespace"], v["token"], v["action"]
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
		err = crd.CreateServiceToken(svc, constant.IdLabel, p.Id)
		if err != nil {
			return makeStringError(err)
		}
	case "delete":
		err = crd.DeleteServiceToken(namespace, token)
		if err != nil {
			return makeStringError(err)
		}
	}
	return nil
}

func (a *App) watchServiceTokens(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("service", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.ServiceTokens(namespace, name)).ServeHTTP(w, r)
}
