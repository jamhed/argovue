package app

import (
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) controlPvcMounts(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, service, action := v["name"], v["namespace"], v["service"], v["action"]
	if err := authObj("pvc", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "create":
		pvc, err := kube.GetPvc(name, namespace)
		if err != nil {
			return makeStringError(err)
		}
		data := struct {
			Owner     string `json:"owner"`
			Canonical string `json:"canonical"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return makeStringError(err)
		}
		label, owner, err := verifyOwner(p, data.Owner)
		if err != nil {
			return makeStringError(err)
		}
		log.Debugf("Deploy mount %s/%s canonical:%s label:%s value:%s", namespace, name, data.Canonical, label, owner)
		err = crd.CreatePvcMount(pvc, a.Args().Namespace(), a.Args().Release(), label, owner, data.Canonical)
		if err != nil {
			return makeStringError(err)
		}
	case "delete":
		err = crd.DeleteInstance(namespace, service)
		if err != nil {
			return makeStringError(err)
		}
	}
	return nil
}

func (a *App) watchPvcMounts(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("pvc", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.PvcMounts(namespace, name)).ServeHTTP(w, r)
}
