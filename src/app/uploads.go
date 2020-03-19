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

func (a *App) controlDatasourceSyncs(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, sync, action := v["name"], v["namespace"], v["sync"], v["action"]
	if err := authObj("datasource", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "create":
		datasource, err := kube.GetDatasource(name, namespace)
		if err != nil {
			return makeStringError(err)
		}
		creds, err := a.Args().AWS().GetCreds(datasource.Spec.Location)
		if err != nil {
			log.Errorf("Error getting session credentials for path:%s, error:%s", datasource.Spec.Location, err)
			return makeStringError(err)
		}
		credsValue, err := creds.Get()
		if err != nil {
			log.Errorf("Error decoding credentials for path:%s, error:%s", datasource.Spec.Location, err)
			return makeStringError(err)
		}
		params := a.Args().RcloneParams()
		params.Key = credsValue.AccessKeyID
		params.Secret = credsValue.SecretAccessKey
		params.Session = credsValue.SessionToken
		err = crd.SyncPvcDatasource(datasource, constant.IdLabel, p.Id, params)
		if err != nil {
			return makeStringError(err)
		}
	case "delete":
		err = crd.DeleteDatasourceSync(namespace, sync)
		if err != nil {
			return makeStringError(err)
		}
	}
	return nil
}

func (a *App) watchDatasourceSyncs(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("datasource", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.DatasourceUploads(namespace, name)).ServeHTTP(w, r)
}

func (a *App) watchJobPods(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("job", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.JobPods(namespace, name)).ServeHTTP(w, r)
}
