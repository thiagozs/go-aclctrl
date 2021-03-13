package config

import (
	"acl-test-go/database"
	"acl-test-go/model"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/thiagozs/go-acl"
)

const (
	aclRead  = "read"
	aclWrite = "write"
	aclList  = "list"
)

// print the contents of the object
func prettyPrint(data interface{}) {
	var p []byte
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", p)
}

func New(ctx context.Context, db database.Repo, secret string) (*acl.ACL, error) {

	// build the models for ACL
	m := acl.NewModel()
	m.Resource("copperbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)
	m.Resource("silverbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)
	m.Resource("diamondbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)
	m.Resource("goldbank").
		Capabilities(aclRead, aclWrite, aclList).
		Alias("read", aclList, aclRead).
		Alias("write", aclWrite, aclRead, aclWrite)

	/*
		// find rules and permission on database
		// Simulate from memory for validation rules
		perm := storage.NewStorage(secret)
		fmt.Printf("[Mock] UserSecrets----> %+v\n", perm.Tokens)
		fmt.Printf("[Mock] UserPolicies---> %+v\n", perm.Policies)
	*/

	secrets, policies, err := db.GetPermissonsBySecret(secret)
	if err != nil {
		log.Fatal(err)
	}
	perm := model.BindDBPermissions(secrets, policies)

	// build config with models ACL
	config := &acl.ResolverConfig{
		Model: m,
		SecretResolver: func(ctx context.Context, s string) (acl.Token, error) {
			return perm.FindTokenBySecret(s)
		},
		PolicyResolver: func(ctx context.Context, p string) (acl.Policy, error) {
			return perm.GetPolicyByName(p)
		},
	}

	resolver, err := acl.NewResolver(config)
	if err != nil {
		return nil, err
	}

	acl, err := resolver.ResolveSecret(ctx, secret)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n\nACL ResolveSecret:\n------------------\n%s\n", acl.String())
	return acl, nil
}
