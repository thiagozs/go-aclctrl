package model

import (
	"acl-test-go/database/gorm_helper"
	"time"
)

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
