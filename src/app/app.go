package app

import (
	"argovue/args"
	"argovue/auth"
	"argovue/authcache"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	redisstore "gopkg.in/boj/redistore.v1"
)

type BrokerMap map[string]map[string]*CrdBroker

type App struct {
	args      *args.Args
	auth      *auth.Auth
	authCache *authcache.Cache
	store     *redisstore.RediStore
	wg        sync.WaitGroup
	brokers   BrokerMap
	config    *crd.Crd
	groups    map[string]string
	subset    BrokerMap
	events    chan *Event
	ver       map[string]interface{}
}

func (a *App) Args() *args.Args {
	return a.args
}

func (a *App) Auth() *auth.Auth {
	return a.auth
}

func (a *App) Store() *redisstore.RediStore {
	return a.store
}

func New() *App {
	a := new(App)
	var err error
	a.args = args.New().LogLevel()
	serializer := mkSerializer(a.Args().SessionKey())
	a.store, err = redisstore.NewRediStore(10, "tcp", a.Args().Redis(), "", serializer.Key())
	if err != nil {
		log.Panicf("Can't create redis session store, error:%s", err)
	}
	a.store.SetSerializer(serializer)
	a.store.Options.Domain = a.Args().BaseDomain()
	a.brokers = make(BrokerMap)
	a.subset = make(BrokerMap)
	a.authCache = authcache.New()
	a.events = make(chan *Event)
	gob.Register(map[string]interface{}{})
	gob.Register(profile.Profile{})
	go a.Serve()
	a.auth = auth.New(a.Args().OIDC())
	a.config = crd.New("argovue.io", "v1", "appconfigs").
		SetFieldSelector(fmt.Sprintf("metadata.namespace=%s,metadata.name=%s-config", a.Args().Namespace(), a.Args().Release())).
		Watch()
	go a.ListenForConfig()
	a.ver = make(map[string]interface{})
	return a
}

func (a *App) SetVersion(version, commit, builddate string) *App {
	a.ver["version"] = version
	a.ver["commit"] = commit
	a.ver["builddate"] = builddate
	if client, err := kube.GetClient(); err == nil {
		if ver, err := client.ServerVersion(); err == nil {
			a.ver["kubernetes"] = ver
		}
	}
	return a
}

var bypassAuth []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("^/profile$"),
	regexp.MustCompile("^/auth"),
	regexp.MustCompile("^/logout$"),
	regexp.MustCompile("^/callback.*$"),
	regexp.MustCompile("^/ui/.*$"),
	regexp.MustCompile("^/proxy/.*$"),
	regexp.MustCompile("^/domain/.*$"),
}

func makeError(code int, format string, args ...interface{}) *appError {
	return &appError{Error: fmt.Sprintf(format, args...), Code: code}
}

func makeStringError(err error) *appError {
	return &appError{Error: fmt.Sprintf("%s", err)}
}

type appError struct {
	Error string
	Code  int
}

type appHandler func(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError
type controlHandler func(p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError
type appErrHandler func(w http.ResponseWriter, r *http.Request) *appError
type httpErrHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn httpErrHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Debug(err.Error)
		http.Error(w, err.Error, err.Code)
	}
}

func (a *App) appHandler(fn appHandler) httpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		session, err := a.Store().Get(r, "auth-session")
		if err != nil {
			return makeError(http.StatusInternalServerError, "Can't get session, error:%s", err)
		}
		pf := session.Values["profile"].(profile.Profile)
		return fn(session.ID, &pf, w, r)
	}
}

func (a *App) appErrHandler(fn appErrHandler) httpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		return fn(w, r)
	}
}

func (a *App) controlHandler(fn controlHandler) httpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		v := mux.Vars(r)
		action := v["action"]
		session, err := a.Store().Get(r, "auth-session")
		if err != nil {
			sendError(w, action, err)
			return nil
		}
		pf, ok := session.Values["profile"].(profile.Profile)
		if !ok {
			sendError(w, action, fmt.Errorf("no profile in session"))
			return nil
		}
		if err := fn(&pf, w, r); err == nil {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
		} else {
			log.Errorf("Action:%s, error:%s", action, err)
			sendError(w, action, fmt.Errorf("%s", err.Error))
		}
		return nil
	}
}

func authObj(kind, name, namespace string, p *profile.Profile) *appError {
	obj, err := kube.GetByKind(kind, name, namespace)
	if err != nil {
		return makeError(http.StatusNotFound, "Can't find object by kind %s/%s/%s, err:%s", kind, namespace, name, err)
	}
	if !p.Authorize(obj) {
		return makeError(http.StatusForbidden, "Not authorized to access object %s/%s/%s", kind, namespace, name)
	}
	return nil
}

func (a *App) authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", a.Args().UIRootDomain())
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for _, re := range bypassAuth {
			if re.MatchString(r.RequestURI) {
				log.Debugf("HTTP: skip auth, from:%s %v", r.RemoteAddr, r.RequestURI)
				next.ServeHTTP(w, r)
				return
			}
		}
		session, err := a.Store().Get(r, "auth-session")
		if err != nil {
			log.Debugf("Can't get session, error:%s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, ok := session.Values["profile"].(profile.Profile)
		if !ok {
			log.Debugf("Not authorized, request:%s", r.RequestURI)
			http.Redirect(w, r, "/auth?redirect="+url.PathEscape(r.RequestURI), http.StatusFound)
			return
		}
		if len(a.brokers[session.ID]) == 0 {
			log.Debugf("Restore brokers for session:%s", session.ID)
			a.onLogin(session.ID, &p)
		}
		log.Debugf("HTTP: %s from:%s %s", p.Id, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func sendError(w http.ResponseWriter, action string, err error) {
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "action": action, "message": fmt.Sprintf("%s", err)})
}

func (a *App) Serve() {
	a.wg.Add(1)
	defer a.wg.Done()
	bindAddr := fmt.Sprintf("%s:%d", a.Args().BindAddr(), a.Args().Port())
	log.Infof("HTTP: at %s static:%s start", bindAddr, a.Args().Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.Args().Dir()))))
	r.HandleFunc("/events", a.handleEvents)
	r.HandleFunc("/profile", a.Profile)
	r.HandleFunc("/objects", a.Objects)
	r.HandleFunc("/version", a.Version)
	r.HandleFunc("/auth", a.AuthInitiate)
	r.Handle("/callback", httpErrHandler(a.AuthCallback))
	r.HandleFunc("/logout", a.Logout)
	r.Handle("/watch/{kind}", a.appHandler(a.watchKind))

	r.Handle("/proxy/{namespace}/{name}/{port}/{rest:.*}", a.appErrHandler(a.proxyService))
	r.Handle("/proxy/{namespace}/{name}/{port}", a.appErrHandler(a.proxyService))
	r.Handle("/domain/{rest:.*}", a.appErrHandler(a.proxyDomain))

	r.Handle("/k8s/{kind}/{namespace}/{name}", a.appHandler(a.watchObject))
	r.Handle("/k8s/pod/{namespace}/{name}/container/{container}/logs", a.appHandler(a.watchPodLogs))
	r.Handle("/k8s/service/{namespace}/{name}/tokens", a.appHandler(a.watchServiceTokens))
	r.Handle("/k8s/service/{namespace}/{name}/tokens/{action}", a.controlHandler(a.controlServiceTokens))
	r.Handle("/k8s/service/{namespace}/{name}/token/{token}/{action}", a.controlHandler(a.controlServiceTokens))
	r.Handle("/k8s/service/{namespace}/{name}/ingresses", a.appHandler(a.watchServiceIngresses))
	r.Handle("/k8s/service/{namespace}/{name}/ingresses/{action}", a.controlHandler(a.controlServiceIngresses))
	r.Handle("/k8s/service/{namespace}/{name}/ingress/{ingress}/{action}", a.controlHandler(a.controlServiceIngresses))

	r.Handle("/k8s/pvc/{namespace}/{name}/datasets", a.appHandler(a.watchPvcDatasets))
	r.Handle("/k8s/pvc/{namespace}/{name}/datasets/{action}", a.controlHandler(a.controlPvcDatasets))
	r.Handle("/k8s/pvc/{namespace}/{name}/dataset/{dataset}/{action}", a.controlHandler(a.controlPvcDatasets))

	r.Handle("/k8s/pvc/{namespace}/{name}/mounts", a.appHandler(a.watchPvcMounts))
	r.Handle("/k8s/pvc/{namespace}/{name}/mounts/{action}", a.controlHandler(a.controlPvcMounts))
	r.Handle("/k8s/pvc/{namespace}/{name}/mount/{service}/{action}", a.controlHandler(a.controlPvcMounts))

	r.Handle("/k8s/dataset/{namespace}/{name}/syncs", a.appHandler(a.watchDatasetSyncs))
	r.Handle("/k8s/dataset/{namespace}/{name}/syncs/{action}", a.controlHandler(a.controlDatasetSyncs))
	r.Handle("/k8s/dataset/{namespace}/{name}/sync/{sync}/{action}", a.controlHandler(a.controlDatasetSyncs))

	r.Handle("/k8s/dataset/{namespace}/{name}/pvcs", a.appHandler(a.watchDatasetPvcs))
	r.Handle("/k8s/dataset/{namespace}/{name}/pvcs/{action}", a.controlHandler(a.controlDatasetPvcs))
	r.Handle("/k8s/dataset/{namespace}/{name}/pvc/{pvc}/{action}", a.controlHandler(a.controlDatasetPvcs))

	r.Handle("/k8s/job/{namespace}/{name}/pods", a.appHandler(a.watchJobPods))

	r.Handle("/catalogue/{namespace}/{name}", a.appHandler(a.watchCatalogue))
	r.Handle("/catalogue/{namespace}/{name}/instances", a.appHandler(a.watchCatalogueInstances))
	r.Handle("/catalogue/{namespace}/{name}/resources", a.appHandler(a.watchCatalogueResources))
	r.Handle("/catalogue/{namespace}/{name}/instance/{instance}", a.appHandler(a.watchCatalogueInstance))
	r.Handle("/catalogue/{namespace}/{name}/instance/{instance}/resources", a.appHandler(a.watchCatalogueInstanceResources))
	r.Handle("/catalogue/{namespace}/{name}/{action}", a.appHandler(a.controlCatalogue)).Methods("POST", "OPTIONS")
	r.Handle("/catalogue/{namespace}/{name}/instance/{instance}/action/{action}", a.appHandler(a.controlCatalogueInstance)).Methods("POST", "OPTIONS")

	r.Handle("/workflow/{namespace}/{name}", a.appHandler(a.watchWorkflow))
	r.Handle("/workflow/{namespace}/{name}/services", a.appHandler(a.watchWorkflowServices))
	r.Handle("/workflow/{namespace}/{name}/mounts", a.appHandler(a.watchWorkflowMounts))
	r.Handle("/workflow/{namespace}/{name}/service/{service}/action/{action}", a.appHandler(a.controlWorkflowService)).Methods("POST", "OPTIONS")
	r.Handle("/workflow/{namespace}/{name}/action/{action}", a.appHandler(a.controlWorkflow)).Methods("POST", "OPTIONS")

	r.Use(a.authMiddleWare)
	srv := &http.Server{
		Handler: r,
		Addr:    bindAddr,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}

func (a *App) Wait() {
	a.wg.Wait()
}
