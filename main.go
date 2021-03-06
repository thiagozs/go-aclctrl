package main

import (
	"acl-test-go/config"
	"acl-test-go/database"
	"acl-test-go/model"
	"acl-test-go/utils"
	"context"
	"fmt"
	"log"

	"github.com/thiagozs/go-acl"
)

func main() {

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	models := []interface{}{
		model.Token{},
		model.Policie{},
		model.Rule{},
	}
	if err := db.CreateTables(models); err != nil {
		log.Fatal(err)
	}

	strp1 := fmt.Sprintf("firstbanks-%s", utils.RandStringRunes(5))
	policy1 := db.CreatePolicy(strp1, []model.Rule{{
		Capabilities: []string{"read"},
		Path:         "*",
		Resource:     fmt.Sprintf("goldbank-%s", utils.RandStringRunes(5)),
	}, {
		Capabilities: []string{"write"},
		Path:         "*",
		Resource:     fmt.Sprintf("silverbank-%s", utils.RandStringRunes(5)),
	}})

	strp2 := fmt.Sprintf("lastbank-%s", utils.RandStringRunes(5))
	policy2 := db.CreatePolicy(strp2, []model.Rule{{
		Capabilities: []string{"read", "write", "list"},
		Path:         "*",
		Resource:     fmt.Sprintf("copperbank-%s", utils.RandStringRunes(5)),
	}})

	token, err := db.CreateSecret(false, 1, []string{policy1.Name, policy2.Name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	perm, policies, err := db.GetPermissonsBySecret(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UserSecret----> %+v\n", perm.Secret)
	fmt.Printf("UserPolicies--> %+v\n", policies)

	config := config.New()
	resolver, err := acl.NewResolver(config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	acl, err := resolver.ResolveSecret(ctx, perm.Secret)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ACL ResolveSecret:\n------------------\n%s\n", acl.String())

	if err := acl.CheckAuthorized(ctx, "silverbank", "*", "write"); err != nil {
		fmt.Println("Not Authorized")
		return
	}
	if err := acl.CheckAuthorized(ctx, "copperbank", "*", "update"); err != nil {
		fmt.Println("Not Authorized")
		return
	}
	fmt.Println("Keep the flow")

}
