package main

import (
	"acl-test-go/config"
	"acl-test-go/database"
	"acl-test-go/model"
	"acl-test-go/utils"
	"context"
	"fmt"
	"log"

	"github.com/edufschmidt/go-acl"
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

	strg := fmt.Sprintf("goldbank-%s", utils.RandStringRunes(5))
	goldBankRules := db.CreateRule(strg, "*", []string{"read"})

	strs := fmt.Sprintf("silverbank-%s", utils.RandStringRunes(5))
	silverBankRules := db.CreateRule(strs, "*", []string{"write"})

	strc := fmt.Sprintf("copperbank-%s", utils.RandStringRunes(5))
	copperBankRules := db.CreateRule(strc, "*", []string{"read", "write", "list"})

	strp1 := fmt.Sprintf("firstbanks-%s", utils.RandStringRunes(5))
	policy1 := db.CreatePolicy(strp1, []model.Rule{goldBankRules, silverBankRules})
	strp2 := fmt.Sprintf("lastbank-%s", utils.RandStringRunes(5))
	policy2 := db.CreatePolicy(strp2, []model.Rule{copperBankRules})

	token := db.CreateToken(false, 1, []string{policy1.Name, policy2.Name})
	fmt.Println(token)

	tokenResponse, policiesResponse, err := db.GetPermissonsByToken(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UserToken    -> %+v\n", tokenResponse.Secret)
	fmt.Printf("UserPolicies -> %+v\n", policiesResponse)

	config := config.New()
	resolver, err := acl.NewResolver(config)
	if err != nil {
		panic(err)
	}

	secret := "54c06ace-7da6-443b-a5a2-05da5294fbd5"
	ctx := context.Background()

	acl, err := resolver.ResolveSecret(ctx, secret)
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
