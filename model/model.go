package model

import (
	"acl-test-go/database/gorm_helper"
	"fmt"
	"time"

	"github.com/thiagozs/go-acl"
)

type PermissionsResponse struct {
	Tokens   []*Token   `json:"tokens,omitempty"`
	Policies []*Policie `json:"policies,omitempty"`
}

type Token struct {
	ID         uint                  `json:"id,omitempty" gorm:"column:id;primaryKey"`
	UserId     uint                  `json:"user_id,omitempty"`
	Privileged bool                  `json:"privileged,omitempty" gorm:"default:false"`
	Secret     string                `json:"secret,omitempty"`
	Policies   gorm_helper.StringArr `json:"policies,omitempty" gorm:"type:text"`
	CreatedAt  time.Time             `json:"created_at,omitempty"`
	UpdatedAt  time.Time             `json:"updated_at,omitempty"`
}

type Policie struct {
	ID        uint      `json:"id,omitempty" gorm:"column:id;primaryKey"`
	Name      string    `json:"name,omitempty"`
	Rules     []Rule    `json:"rules,omitempty" gorm:"foreignKey:PoliceID"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Rule struct {
	ID           uint                  `json:"id,omitempty" gorm:"column:id;primaryKey"`
	PoliceID     uint                  `json:"police_id,omitempty"`
	Resource     string                `json:"resource,omitempty"`
	Path         string                `json:"path,omitempty"`
	Capabilities gorm_helper.StringArr `json:"capabilities,omitempty" gorm:"type:text"`
	CreatedAt    time.Time             `json:"created_at,omitempty" `
	UpdatedAt    time.Time             `json:"updated_at,omitempty"`
}

func (t Token) TableName() string {
	return "acl_tokens"
}

func (p Policie) TableName() string {
	return "acl_policies"
}

func (r Rule) TableName() string {
	return "acl_rules"
}

func (r *PermissionsResponse) FindTokenBySecret(s string) (acl.Token, error) {
	for _, t := range r.Tokens {
		if t.Secret == s {
			return t, nil
		}
	}
	return nil, nil
}

func (r *PermissionsResponse) GetPolicyByName(n string) (acl.Policy, error) {
	for _, p := range r.Policies {
		if p.Name == n {
			return p, nil
		}
	}
	return nil, fmt.Errorf("not found : %s", n)
}

func (t *Token) PermPolicies() []string {
	return t.Policies
}

func (t *Token) PermIsPrivileged() bool {
	return t.Privileged
}

func (p *Policie) PermName() string {
	return p.Name
}

func (p *Policie) PermRules() []acl.Rule {
	rules := []acl.Rule{}
	for _, r := range p.Rules {
		rules = append(rules, r)
	}
	return rules
}

func (r Rule) GetResource() string {
	return r.Resource
}

func (r Rule) GetPath() string {
	return r.Path
}

func (r Rule) GetCapabilities() []string {
	return r.Capabilities
}
