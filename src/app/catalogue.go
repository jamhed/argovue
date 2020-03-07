package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"encoding/json"
	"fmt"
	"net/http"

	argovuev1 "argovue/apis/argovue.io/v1"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (a *App) watchCatalogue(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("argovue", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.Catalogue(namespace, name)).ServeHTTP(w, r)
}

func (a *App) watchCatalogueInstances(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("argovue", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.CatalogueInstances(namespace, name)).ServeHTTP(w, r)
}

func (a *App) watchCatalogueResources(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("argovue", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.CatalogueResources(namespace, name)).ServeHTTP(w, r)
}

func (a *App) watchCatalogueInstanceResources(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, instance := v["name"], v["namespace"], v["instance"]
	if err := authObj("helmrelease", instance, namespace, p); err != nil {
		return err
	}
	id := fmt.Sprintf("%s-%s-%s-resources", namespace, name, instance)
	cb := a.maybeNewIdSubsetBroker(session.ID, id)
	cb.AddCrd(crd.CatalogueWorkflows(namespace, instance))
	cb.AddCrd(crd.CatalogueInstancePods(namespace, instance))
	cb.AddCrd(crd.CatalogueInstancePvcs(namespace, instance))
	cb.AddCrd(crd.CatalogueInstanceServices(namespace, instance))
	cb.AddCrd(crd.CatalogueInstanceIngresses(namespace, instance))
	return cb.ServeHTTP(w, r)
}

func (a *App) watchCatalogueInstance(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, instance := v["name"], v["namespace"], v["instance"]
	if err := authObj("helmrelease", instance, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.CatalogueInstance(namespace, name, instance)).ServeHTTP(w, r)
}

func (a *App) controlCatalogue(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, action := v["name"], v["namespace"], v["action"]
	if err := authObj("argovue", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "deploy":
		var svc *argovuev1.Service
		var label, owner string

		svc, err = kube.GetArgovueService(name, namespace)
		if err != nil {
			goto err
		}
		data := struct {
			Owner     string                 `json:"owner"`
			Canonical string                 `json:"canonical"`
			Input     []argovuev1.InputValue `json:"input"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			goto err
		}
		label, owner, err = verifyOwner(p, data.Owner)
		if err != nil {
			goto err
		}
		log.Debugf("Deploy service %s/%s canonical:%s label:%s value:%s", namespace, name, data.Canonical, label, owner)
		err = crd.Deploy(svc, label, owner, data.Canonical, data.Input)
		if err != nil {
			goto err
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	return nil
err:
	log.Errorf("Can't %s catalogue %s/%s, error:%s", action, namespace, name, err)
	sendError(w, action, err)
	return nil
}

func (a *App) controlCatalogueInstance(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, instance, action := v["name"], v["namespace"], v["instance"], v["action"]
	if err := authObj("helmrelease", instance, namespace, p); err != nil {
		return err
	}
	switch action {
	case "delete":
		err := crd.DeleteInstance(namespace, instance)
		if err != nil {
			log.Errorf("Can't %s catalogue %s/%s instance:%s, error:%s", action, namespace, name, instance, err)
			sendError(w, action, err)
		} else {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
		}
	}
	return nil
}

func verifyOwner(p *profile.Profile, owner string) (string, string, error) {
	if p.Id == owner {
		return constant.IdLabel, p.Id, nil
	}
	for _, g := range p.EffectiveGroups {
		if g == owner {
			return constant.GroupLabel, owner, nil
		}
	}
	return "", "", fmt.Errorf("Can't verify owner:%s", owner)
}
