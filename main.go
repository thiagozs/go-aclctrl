package main

import (
	"acl-test-go/config"
	"acl-test-go/database"
	"acl-test-go/model"
	"acl-test-go/utils"
	"context"
	"fmt"
	"log"
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

	strp1 := fmt.Sprintf("firstbanks%s", utils.RandStringRunes(5))
	policy1 := db.CreatePolicy(strp1, 1, []model.Rule{{
		Capabilities: []string{"read"},
		Path:         "*",
		Resource:     "goldbank",
	}, {
		Capabilities: []string{"write"},
		Path:         "*",
		Resource:     "silverbank",
	}})

	strp2 := fmt.Sprintf("lastbank%s", utils.RandStringRunes(5))
	policy2 := db.CreatePolicy(strp2, 1, []model.Rule{{
		Resource:     "copperbank",
		Path:         "*",
		Capabilities: []string{"read", "write", "list"},
	}})

	secret, err := db.CreateSecret(false, 1, []string{policy1.Name, policy2.Name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secret)

	ctx := context.Background()

	acl, err := config.New(ctx, db, secret)
	if err != nil {
		log.Fatal(err)
	}

	if err := acl.CheckAuthorized(ctx, "goldbank", "*", "list"); err != nil {
		fmt.Println("Not Authorized")
		return
	}
	if err := acl.CheckAuthorized(ctx, "silverbank", "*", "write"); err != nil {
		fmt.Println("Not Authorized")
		return
	}
	fmt.Println("Keep the flow")

}
