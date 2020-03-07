package app

import (
	"argovue/crd"
	"argovue/profile"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) watchPodLogs(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, container, kind := v["name"], v["namespace"], v["container"], "pod"
	if err := authObj(kind, name, namespace, p); err != nil {
		return err
	}
	return a.streamPodLogs(w, r, name, namespace, container)
}

func (a *App) watchObject(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	kind, name, namespace := v["kind"], v["name"], v["namespace"]
	if err := authObj(kind, name, namespace, p); err != nil {
		return err
	}
	crd, err := crd.GetByKind(kind, namespace, name)
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't create watcher %s/%s/%s", kind, namespace, name)
	}
	return a.maybeNewSubsetBroker(sid, crd).ServeHTTP(w, r)
}

func (a *App) watchKind(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	kind := vars["kind"]
	broker := a.getBroker(sid, kind)
	if broker == nil {
		return makeError(http.StatusNotFound, "Can't find broker by kind:%s", kind)
	}
	return broker.ServeHTTP(w, r)
}
