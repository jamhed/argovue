package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/profile"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type BrokerDef struct {
	Name     string
	Group    string
	Version  string
	Resource string
}

func (a *App) onLogout(sessionId string) {
	sessionData, ok := a.brokers[sessionId]
	if !ok {
		return
	}
	for name, _ := range sessionData {
		log.Debugf("Delete broker: %s", name)
		if broker, ok := sessionData[name]; ok {
			broker.Stop()
			delete(sessionData, name)
		}
	}
}

func (a *App) onLogin(sessionId string, p *profile.Profile) {
	brokerDefs := []BrokerDef{
		{"catalogue", "argovue.io", "v1", "services"},
		{"datasources", "argovue.io", "v1", "datasources"},
		{"pvcs", "", "v1", "persistentvolumeclaims"},
		{"workflows", "argoproj.io", "v1alpha1", "workflows"},
		{"services", "", "v1", "services"},
	}
	groupSelector := fmt.Sprintf("%s in (%s)", constant.GroupLabel, strings.Join(p.EffectiveGroups, ","))
	selector := fmt.Sprintf("%s in (%s)", constant.IdLabel, p.IdLabel())
	for _, def := range brokerDefs {
		broker := a.newBroker(sessionId, def.Name)
		if len(p.EffectiveGroups) > 0 {
			broker.AddCrd(crd.New(def.Group, def.Version, def.Resource).SetLabelSelector(groupSelector))
		}
		broker.AddCrd(crd.New(def.Group, def.Version, def.Resource).SetLabelSelector(selector))
	}
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Redirect(w, r, a.Args().UIRootURL(), http.StatusFound)
		return
	}
	a.onLogout(session.ID)
	delete(session.Values, "state")
	delete(session.Values, "auth-session")
	delete(session.Values, "profile")
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		log.Errorf("Can't delete session, error:%s", err)
	}
	http.Redirect(w, r, a.Args().UIRootURL(), http.StatusFound)
}

func (a *App) AuthInitiate(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirect := r.URL.Query().Get("redirect")

	state := base64.StdEncoding.EncodeToString(b)
	session, _ := a.Store().Get(r, "auth-session")
	session.Values["state"] = state
	if len(redirect) > 0 {
		unescape, err := url.PathUnescape(redirect)
		if err == nil {
			log.Debugf("AUTH: keep redirect value:%s", unescape)
			session.Values["redirect"] = unescape
		} else {
			log.Debugf("AUTH: error unescape path:%s, error:%s", redirect, err)
		}
	}

	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, a.Auth().Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (a *App) AuthCallback(w http.ResponseWriter, r *http.Request) *appError {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't get session, error:%s", err.Error())
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		return makeError(http.StatusBadRequest, "Can't get proper state value from request")
	}

	token, err := a.Auth().Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		return makeError(http.StatusUnauthorized, "Can't get token by code value")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return makeError(http.StatusInternalServerError, "Can't get id_token value from token:%s", token)
	}

	idToken, err := a.Auth().Provider.Verifier(&oidc.Config{ClientID: a.Auth().Config.ClientID}).Verify(context.TODO(), rawIDToken)
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't verify token, error:%s", err)
	}

	var idTokenClaims map[string]interface{}
	if err := idToken.Claims(&idTokenClaims); err != nil {
		return makeError(http.StatusInternalServerError, "Can't decode token claims, error:%s", err)
	}

	log.Debugf("OIDC: id token claims: %s", idTokenClaims)
	p := profile.New().FromMap(idTokenClaims, a.Args().OidcUserId())
	if len(p.Groups) == 0 {
		if userInfoclaims, err := a.userInfo(token); err == nil {
			log.Debugf("OIDC: user info claims: %s", userInfoclaims)
			p.FromMap(userInfoclaims, a.Args().OidcUserId())
		} else {
			log.Debugf("Can't get userInfo claims, error:%s", err)
		}
	}
	p.MapValues(a.groups)

	session.Values["profile"] = p
	a.onLogin(session.ID, p)

	redirectUrl := session.Values["redirect"]
	delete(session.Values, "redirect")

	if err = session.Save(r, w); err != nil {
		return makeError(http.StatusInternalServerError, "Can't save session, error:%s", err)
	}
	// this is to set cookie for the api domain
	redirect := `<html><head><script type="text/javascript">window.location.href="%s"</script></head><body></body></html>`
	if re, ok := redirectUrl.(string); ok {
		log.Debugf("AUTH: redirecting to:%s", re)
		fmt.Fprintf(w, redirect, re)
	} else {
		fmt.Fprintf(w, redirect, a.Args().UIRootURL())
	}
	return nil
}

func (a *App) userInfo(token *oauth2.Token) (map[string]interface{}, error) {
	var claims map[string]interface{}
	userinfo, err := a.Auth().Provider.UserInfo(context.TODO(), a.Auth().Config.TokenSource(context.TODO(), token))
	if err != nil {
		log.Errorf("Can't request user info, error:%s", err)
		return nil, err
	}
	err = userinfo.Claims(&claims)
	if err != nil {
		log.Errorf("Can't decode claims, error:%s", err)
		return nil, err
	}
	return claims, nil
}
