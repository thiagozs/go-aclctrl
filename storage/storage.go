package storage

import (
	"fmt"

	"github.com/thiagozs/go-acl"
)

type storage struct {
	Tokens   []*Tokens
	Policies []*Policies
}

func NewStorage(secret string) *storage {
	return &storage{
		Tokens: []*Tokens{
			//	{true, "39076595-19a6-4582-b0d9-bb4a266fd48a", []string{"policy1"}},
			//	{false, "71036287-81d1-4748a-b4d5-25c2ee6f57ae", []string{"policy2"}},
			{false, secret, []string{"policy3", "policy4"}},
			//	{false, "e690413b-827b-400e-bc38-92a4b1580eac", []string{"policy1", "policy2"}},
		},
		Policies: []*Policies{
			// {"policy1", []Rules{
			// 	{"silverbank", "*", []string{"read"}},
			// 	{"goldbank", "*", []string{"read"}},
			// 	{"copperbank", "", []string{"read"}},
			// }},
			// {"policy2", []Rules{
			// 	{"silverbank", "*", []string{"write"}},
			// 	{"goldbank", "*", []string{"deny"}},
			// 	{"copperbank", "*", []string{"list"}},
			// }},
			{"policy3", []*Rules{
				{"silverbank", "pol3", []string{"read"}},
				{"copperbank", "pol4", []string{"write"}},
			}},
			{"policy4", []*Rules{
				{"goldbank", "pol1", []string{"read", "write", "list"}},
				{"silverbank", "pol2", []string{"write"}},
			}},
			// {"anonymous", []Rules{
			// 	{"copperbank", "*", []string{"list"}},
			// 	{"silverbank", "*", []string{"list"}},
			// }},
		},
	}
}

func (r *storage) FindTokenBySecret(s string) (acl.Token, error) {
	for _, t := range r.Tokens {
		if t.Secret == s {
			return t, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", s)
}

func (r *storage) GetPolicyByName(n string) (acl.Policy, error) {
	for _, p := range r.Policies {
		if p.Name == n {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", n)
}

type Tokens struct {
	Privileged bool
	Secret     string
	Policies   []string
}

func (t *Tokens) PermPolicies() []string {
	return t.Policies
}

func (t *Tokens) PermIsPrivileged() bool {
	return t.Privileged
}

type Policies struct {
	Name  string
	Rules []*Rules
}

func (p *Policies) PermName() string {
	return p.Name
}

func (p *Policies) PermRules() []acl.Rule {
	rules := []acl.Rule{}
	for _, r := range p.Rules {
		rules = append(rules, r)
	}
	return rules
}

type Rules struct {
	Resource     string
	Path         string
	Capabilities []string
}

func (r *Rules) GetResource() string {
	fmt.Printf("Resource => %+v\n", r.Resource)
	return r.Resource
}

func (r *Rules) GetPath() string {
	fmt.Printf("Path => %+v\n", r.Path)
	return r.Path
}

func (r *Rules) GetCapabilities() []string {
	fmt.Printf("Capabilities => %+v\n", r.Capabilities)
	return r.Capabilities
}
