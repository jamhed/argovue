package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) controlPvcDatasources(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, datasource, action := v["name"], v["namespace"], v["datasource"], v["action"]
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
		err = crd.CreatePvcDatasource(pvc, constant.IdLabel, p.Id)
		if err != nil {
			return makeStringError(err)
		}
	case "delete":
		err = crd.DeletePvcDatasource(namespace, datasource)
		if err != nil {
			return makeStringError(err)
		}
	}
	return nil
}

func (a *App) watchPvcDatasources(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("pvc", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.PvcDatasources(namespace, name)).ServeHTTP(w, r)
}

func (a *App) controlDatasourcePvcs(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, pvcName, action := v["name"], v["namespace"], v["pvc"], v["action"]
	if err := authObj("datasource", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "create":
		ds, err := kube.GetDatasource(name, namespace)
		if err != nil {
			return makeStringError(err)
		}
		err = crd.CreateDatasourcePvc(ds, constant.IdLabel, p.Id)
		if err != nil {
			return makeStringError(err)
		}
	case "sync":
		ds, err := kube.GetDatasource(name, namespace)
		if err != nil {
			return makeStringError(err)
		}
		pvc, err := kube.GetPvc(pvcName, namespace)
		if err != nil {
			return makeStringError(err)
		}
		creds, err := a.Args().AWS().GetCreds(ds.Spec.Location)
		if err != nil {
			log.Errorf("Error getting session credentials for path:%s, error:%s", ds.Spec.Location, err)
			return makeStringError(err)
		}
		credsValue, err := creds.Get()
		if err != nil {
			log.Errorf("Error decoding credentials for path:%s, error:%s", ds.Spec.Location, err)
			return makeStringError(err)
		}
		params := a.Args().RcloneParams()
		params.Key = credsValue.AccessKeyID
		params.Secret = credsValue.SecretAccessKey
		params.Session = credsValue.SessionToken
		err = crd.SyncDatasourcePvc(ds, pvc, constant.IdLabel, p.Id, params)
		if err != nil {
			return makeStringError(err)
		}
	case "delete":
		err = crd.DeleteDatasourcePvc(namespace, pvcName)
		if err != nil {
			return makeStringError(err)
		}
	}
	return nil
}

func (a *App) watchDatasourcePvcs(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("datasource", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.DatasourcePvcs(namespace, name)).ServeHTTP(w, r)
}
