package config

import (
	"acl-test-go/storage"
	"context"

	"github.com/edufschmidt/go-acl"
)

const (
	aclRead  = "read"
	aclWrite = "write"
	aclList  = "list"
)

func New() *acl.ResolverConfig {

	// build the models for ACL
	model := acl.NewModel()
	model.Resource("copperbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)
	model.Resource("silverbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)
	model.Resource("diamondbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)
	model.Resource("goldbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)

	// find rules and permission on database
	// Simulate from memory for validation rules
	// TODO: change for databaseRepo
	storage := storage.NewStorage()

	// build config with models ACL
	return &acl.ResolverConfig{
		Logger: nil,
		Model:  model,
		SecretResolver: func(ctx context.Context, s string) (acl.Token, error) {
			return storage.FindTokenBySecret(s)
		},
		PolicyResolver: func(ctx context.Context, p string) (acl.Policy, error) {
			return storage.GetPolicyByName(p)
		},
	}
}
