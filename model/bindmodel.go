package model

import (
	"fmt"

	"github.com/thiagozs/go-acl"
)

type PermissionsResponse struct {
	Tokens   []*TokenResponse
	Policies []*PolicieResponse
}

type TokenResponse struct {
	Privileged bool
	Secret     string
	Policies   []string
}

type PolicieResponse struct {
	Name  string
	Rules []*RuleResponse
}

type RuleResponse struct {
	Resource     string
	Path         string
	Capabilities []string
}

func BindDBPermissions(secrets []Token, policies []Policie) *PermissionsResponse {
	resp := &PermissionsResponse{
		Policies: []*PolicieResponse{},
		Tokens:   []*TokenResponse{},
	}

	for _, v := range policies {

		pol := &PolicieResponse{
			Name:  v.Name,
			Rules: []*RuleResponse{},
		}

		for _, vv := range v.Rules {
			pol.Rules = append(pol.Rules, &RuleResponse{
				Resource:     vv.Resource,
				Path:         vv.Path,
				Capabilities: vv.Capabilities.ToStringArr(),
			})
		}

		resp.Policies = append(resp.Policies, pol)
	}

	for _, v := range secrets {
		resp.Tokens = append(resp.Tokens, &TokenResponse{
			Privileged: v.Privileged,
			Secret:     v.Secret,
			Policies:   v.Policies.ToStringArr(),
		})
	}

	return resp
}

func (r *PermissionsResponse) FindTokenBySecret(s string) (acl.Token, error) {
	for _, t := range r.Tokens {
		if t.Secret == s {
			return t, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", s)
}

func (r *PermissionsResponse) GetPolicyByName(n string) (acl.Policy, error) {
	for _, p := range r.Policies {
		if p.Name == n {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", n)
}

func (r *PermissionsResponse) GetSecrets() []*TokenResponse {
	return r.Tokens
}

func (r *PermissionsResponse) GetPolicies() []*PolicieResponse {
	return r.Policies
}

func (t *TokenResponse) PermPolicies() []string {
	return t.Policies
}

func (t *TokenResponse) PermIsPrivileged() bool {
	return t.Privileged
}

func (p *PolicieResponse) PermName() string {
	return p.Name
}

func (p *PolicieResponse) PermRules() []acl.Rule {
	rules := []acl.Rule{}
	for _, r := range p.Rules {
		rules = append(rules, r)
	}
	return rules
}

func (r *RuleResponse) GetResource() string {
	return r.Resource
}

func (r *RuleResponse) GetPath() string {
	return r.Path
}

func (r *RuleResponse) GetCapabilities() []string {
	return r.Capabilities
}
