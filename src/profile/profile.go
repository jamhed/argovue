package profile

import (
	"argovue/constant"
	"argovue/util"
	"crypto/sha1"
	"encoding/hex"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Profile struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	PreferredName   string   `json:"preferred_username"`
	Email           string   `json:"email"`
	ZoneInfo        string   `json:"zoneinfo"`
	Groups          []string `json:"groups"`
	EffectiveId     string   `json:"effective_id"`
	EffectiveGroups []string `json:"effective_groups"`
	Subject         string   `json:"sub"`
}

func New() *Profile {
	return &Profile{}
}

func (p *Profile) IdLabel() string {
	return GetHash(p.Id)
}

func GetHash(text string) string {
	hasher := sha1.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func MaybeHash(label, owner string) string {
	if label == constant.IdLabel {
		return GetHash(owner)
	}
	return owner
}

func (p *Profile) FromMap(m map[string]interface{}, idKey string) *Profile {
	if s, ok := m["email"].(string); ok {
		p.Email = s
	}
	if s, ok := m["sub"].(string); ok {
		p.Subject = s
	}
	if s, ok := m["name"].(string); ok {
		p.Name = s
	}
	if s, ok := m["preferred_username"].(string); ok {
		p.PreferredName = s
	}
	if s, ok := m["zoneinfo"].(string); ok {
		p.ZoneInfo = s
	}
	if s, ok := m[idKey].(string); ok {
		p.Id = s
	}
	p.EffectiveId = GetHash(p.Id)
	p.Groups = util.Li2s(m["groups"])
	return p
}

func (p *Profile) MapValues(m map[string]string) *Profile {
	effGroups := []string{}
	for _, group := range p.Groups {
		if k8sGroup, ok := m[group]; ok {
			effGroups = append(effGroups, k8sGroup)
		}
	}
	p.EffectiveGroups = effGroups
	if len(p.Id) == 0 {
		p.Id = p.Subject
	}
	p.EffectiveId = GetHash(p.Id)
	return p
}

func (p *Profile) Authorize(obj metav1.Object) bool {
	labels := obj.GetLabels()
	if groupLabel, ok := labels[constant.GroupLabel]; ok && len(groupLabel) > 0 {
		for _, group := range p.EffectiveGroups {
			if group == groupLabel {
				return true
			}
		}
	}
	if idLabel, ok := labels[constant.IdLabel]; ok && len(idLabel) > 0 && idLabel == p.IdLabel() {
		return true
	}
	return false
}
