package app

import (
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/argoproj/argo/workflow/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

func (a *App) watchWorkflow(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.Workflow(namespace, name)).ServeHTTP(w, r)
}

func (a *App) watchWorkflowMounts(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	id := fmt.Sprintf("%s-%s-mounts", namespace, name)
	return a.maybeNewIdSubsetBroker(sid, id).AddCrd(crd.WorkflowMounts(namespace, name)).ServeHTTP(w, r)
}

func (a *App) watchWorkflowServices(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	return a.maybeNewSubsetBroker(sid, crd.WorkflowServices(namespace, name)).ServeHTTP(w, r)
}

func (a *App) controlWorkflowService(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, service, action := v["name"], v["namespace"], v["service"], v["action"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "delete":
		err = crd.DeleteInstance(namespace, service)
	default:
		err = fmt.Errorf("Unknown action:%s", action)
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
	return nil
}

func (a *App) controlWorkflow(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, action := v["name"], v["namespace"], v["action"]

	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}

	wfClientset, err := kube.GetWfClientset()
	if err != nil {
		log.Errorf("Can't get argo clientset, error:%s", err)
		sendError(w, action, err)
		return nil
	}
	wfClient := kube.GetWfClient(wfClientset, namespace)
	wf, err := wfClient.Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get workflow %s/%s, error:%s", namespace, name, err)
		sendError(w, action, err)
		return nil
	}

	kubeClient, _ := kube.GetClient()
	switch action {
	case "retry":
		_, err = util.RetryWorkflow(kubeClient, wfClient, wf)
	case "resubmit":
		newWF, err := util.FormulateResubmitWorkflow(wf, false)
		if err == nil {
			_, err = util.SubmitWorkflow(wfClient, wfClientset, namespace, newWF, nil)
		}
	case "delete":
		err = wfClient.Delete(name, &metav1.DeleteOptions{})
	case "suspend":
		err = util.SuspendWorkflow(wfClient, name)
	case "resume":
		err = util.ResumeWorkflow(wfClient, name)
	case "terminate":
		err = util.TerminateWorkflow(wfClient, name)
	case "mount":
		err = crd.DeployFilebrowser(wf, a.Args().Namespace(), a.Args().Release(), p.Id)
	default:
		err = fmt.Errorf("unrecognized command %s", action)
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
	return nil
}
