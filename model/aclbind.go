package model

type PermissionsResponse struct {
	Tokens   []*TokenResponse   `json:"tokens,omitempty"`
	Policies []*PolicieResponse `json:"policies,omitempty"`
}

type TokenResponse struct {
	Privileged bool     `json:"privileged,omitempty"`
	Secret     string   `json:"secret,omitempty"`
	Policies   []string `json:"policies,omitempty"`
}

type RuleResponse struct {
	Resource     string   `json:"resource,omitempty"`
	Path         string   `json:"path,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}

type PolicieResponse struct {
	Name  string         `json:"name,omitempty"`
	Rules []RuleResponse `json:"rules,omitempty"`
}
