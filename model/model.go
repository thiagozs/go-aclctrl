package model

import (
	"time"

	"gorm.io/datatypes"
)

/* Structs for Response only ----------- */

type PermissionsResponse struct {
	Tokens   []*TokenResponse   `json:"tokens,omitempty"`
	Policies []*PolicieResponse `json:"policies,omitempty"`
}

type TokenResponse struct {
	Id         uint      `json:"id,omitempty" gorm:"primaryKey"`
	UserId     uint      `json:"user_id,omitempty"`
	Privileged bool      `json:"privileged,omitempty" gorm:"default:false"`
	Secret     string    `json:"secret,omitempty"`
	Policies   []string  `json:"policies,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type RuleResponse struct {
	Id           uint      `json:"id,omitempty" gorm:"primaryKey"`
	Resource     string    `json:"resource,omitempty"`
	Path         string    `json:"path,omitempty"`
	Capabilities []string  `json:"capabilities,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" `
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type PolicieResponse struct {
	Id        uint           `json:"id,omitempty" gorm:"primaryKey"`
	Name      string         `json:"name,omitempty"`
	Rules     []RuleResponse `json:"rules,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
}

/* Structs for DB methods ----------- */

type Token struct {
	Id         uint           `json:"id,omitempty" gorm:"primaryKey"`
	UserId     uint           `json:"user_id,omitempty"`
	Privileged bool           `json:"privileged,omitempty" gorm:"default:false"`
	Secret     string         `json:"secret,omitempty"`
	Policies   datatypes.JSON `json:"policies,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
}

type Policie struct {
	Id        uint      `json:"id,omitempty" gorm:"primaryKey"`
	Name      string    `json:"name,omitempty"`
	Rules     []Rule    `json:"rules,omitempty" gorm:"foreignKey:PoliceId"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Rule struct {
	Id           uint           `json:"id,omitempty" gorm:"primaryKey"`
	PoliceId     uint           `json:"police_id,omitempty"`
	Resource     string         `json:"resource,omitempty"`
	Path         string         `json:"path,omitempty"`
	Capabilities datatypes.JSON `json:"capabilities,omitempty"`
	CreatedAt    time.Time      `json:"created_at,omitempty" `
	UpdatedAt    time.Time      `json:"updated_at,omitempty"`
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
