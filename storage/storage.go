package storage

import (
	"fmt"

	"github.com/thiagozs/go-acl"
)

type storage struct {
	tokens   []*Tokens
	policies []*Policies
}

func NewStorage() *storage {
	return &storage{
		tokens: []*Tokens{
			{true, "39076595-19a6-4582-b0d9-bb4a266fd48a", []string{"policy3"}},
			{false, "71036287-81d1-4748a-b4d5-25c2ee6f57ae", []string{"policy1"}},
			{false, "54c06ace-7da6-443b-a5a2-05da5294fbd5", []string{"policy2", "policy4"}},
			{false, "e690413b-827b-400e-bc38-92a4b1580eac", []string{"policy1", "policy2"}},
		},
		policies: []*Policies{
			{"policy1", []*Rules{
				{"silverbank", "*", []string{"read"}},
				{"goldbank", "*", []string{"read"}},
				{"copperbank", "", []string{"read"}},
			}},
			{"policy2", []*Rules{
				{"silverbank", "*", []string{"write"}},
				{"goldbank", "*", []string{"deny"}},
				{"copperbank", "*", []string{"list"}},
			}},
			{"policy3", []*Rules{
				{"silverbank", "*", []string{"read"}},
				{"copperbank", "*", []string{"write"}},
			}},
			{"policy4", []*Rules{
				{"goldbank", "*", []string{"read", "write", "list"}},
				{"silverbank", "*", []string{"write"}},
			}},
			{"anonymous", []*Rules{
				{"copperbank", "*", []string{"list"}},
				{"silverbank", "*", []string{"list"}},
			}},
		},
	}
}

func (r *storage) FindTokenBySecret(s string) (acl.Token, error) {
	for _, t := range r.tokens {
		if t.secret == s {
			return t, nil
		}
	}
	return nil, nil
}

func (r *storage) GetPolicyByName(n string) (acl.Policy, error) {
	for _, p := range r.policies {
		if p.name == n {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", n)
}

type Tokens struct {
	privileged bool
	secret     string
	policies   []string
}

func (t *Tokens) PermPolicies() []string {
	return t.policies
}

func (t *Tokens) PermIsPrivileged() bool {
	return t.privileged
}

type Policies struct {
	name  string
	rules []*Rules
}

func (p *Policies) PermName() string {
	return p.name
}

func (p *Policies) PermRules() []acl.Rule {
	rules := []acl.Rule{}
	for _, r := range p.rules {
		rules = append(rules, r)
	}
	return rules
}

type Rules struct {
	resource     string
	path         string
	capabilities []string
}

func (r *Rules) GetResource() string {
	return r.resource
}

func (r *Rules) GetPath() string {
	return r.path
}

func (r *Rules) GetCapabilities() []string {
	return r.capabilities
}
