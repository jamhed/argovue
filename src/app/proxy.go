package app

import (
	"argovue/kube"
	"argovue/profile"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func serviceKey(auth, namespace, name string) string {
	return fmt.Sprintf("%s-%s-%s", namespace, name, auth)
}

func maybeGetBearer(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		return ""
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		log.Warnf("Invalid authorization header value:%s", reqToken)
		return ""
	}
	return splitToken[1]
}

func authorizeToken(svc *v1.Service, token string) bool {
	tokens, err := kube.GetServiceTokens(svc.Name, svc.Namespace)
	if err != nil {
		log.Debugf("Can't get tokens for service %s/%s, error:%s", svc.Namespace, svc.Name, err)
	}
	for _, t := range tokens.Items {
		if token == t.Spec.Value {
			log.Debugf("Proxy: authenticated by bearer token")
			return true
		}
	}
	return false
}

func (a *App) proxyDomain(w http.ResponseWriter, r *http.Request) *appError {
	re := regexp.MustCompile(`^(.+?)\.(.+?)\.(.+?)\.svc\.cluster`)
	tmp := re.FindAllStringSubmatch(r.Host, -1)
	if !(len(tmp) == 1 && len(tmp[0]) == 4) {
		return makeError(http.StatusForbidden, "Proxy: can't parse host:%s", r.Host)
	}
	vars := mux.Vars(r)
	r = mux.SetURLVars(r, map[string]string{
		"namespace": tmp[0][1],
		"name":      tmp[0][2],
		"port":      tmp[0][3],
		"rest":      vars["rest"],
		"keep":      "true",
	})
	return a.proxyService(w, r)
}

func (a *App) proxyService(w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, port, rest, keep := v["name"], v["namespace"], v["port"], v["rest"], v["keep"]
	var pf profile.Profile
	var ok bool
	var session *sessions.Session
	var svc *v1.Service
	var err error
	var token = maybeGetBearer(r)
	var banner string

	// verify service exists
	svc, err = kube.GetService(name, namespace)
	if err != nil {
		return makeError(http.StatusForbidden, "Proxy: no service %s/%s, access denied, error:%s", namespace, name, err)
	}

	// authenticate by token
	if len(token) > 0 {
		if v, ok := a.authCache.Check(serviceKey(token, namespace, name)); ok {
			banner = v
			goto auth
		}
		if authorizeToken(svc, token) {
			a.authCache.Add(serviceKey(token, namespace, name), token)
			banner = token
			goto auth
		} else {
			return makeError(http.StatusForbidden, "Proxy: %s/%s, no bearer auth, access denied", namespace, name)
		}
	}

	session, err = a.Store().Get(r, "auth-session")
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't get session, error:%s", err)
	}

	if v, ok := a.authCache.Check(serviceKey(session.ID, namespace, name)); ok {
		banner = v
		goto auth
	}

	// authenticate by service ownership
	if pf, ok = session.Values["profile"].(profile.Profile); ok {
		if !pf.Authorize(svc) {
			return makeError(http.StatusForbidden, "Proxy: %s/%s, no auth, access denied", namespace, name)
		} else {
			a.authCache.Add(serviceKey(session.ID, namespace, name), pf.Id)
			banner = pf.Id
			goto auth
		}
	} else {
		fwd := fmt.Sprintf("%s://%s/%s", r.Header.Get("X-Forwarded-Proto"), r.Header.Get("X-Forwarded-Host"), rest)
		redirect := fmt.Sprintf("https://%s/auth?redirect=%s", a.Args().BaseDomain(), url.PathEscape(fwd))
		log.Debugf("Proxy: %s/%s, no profile, access denied, redirecting to:%s", namespace, name, redirect)
		http.Redirect(w, r, redirect, http.StatusFound)
		return nil
	}

auth:
	schema := "http"
	if port == "443" {
		schema = "https"
	}

	target := fmt.Sprintf("%s://%s.%s.svc.cluster.local:%s", schema, name, namespace, port)

	if keep != "true" {
		rest = fmt.Sprintf("/proxy/%s/%s/%s/%s", namespace, name, port, rest)
	}

	log.Debugf("Proxy: %s from:%s target:%s url:%s", banner, r.RemoteAddr, target, rest)

	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.URL.Path = rest
	r.Header.Set("Host", r.Host)
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Forwarded-Proto", "https")
	proxy.ServeHTTP(w, r)
	return nil
}
