# Golang ACL Control - with SQL Database

## Stack of development

* golang 1.15
* Gorm Lib OCR
* ACL Lib
* MySQL Database
* Docker compose and containers

Example code `main.go`
```go
    // Register on database the user with the policies you need
	secret, err := db.CreateSecret(false, 1, []string{policy1.Name, policy2.Name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secret)

	ctx := context.Background()

    // start config with yours custom settins of rules
	acl, err := config.New(ctx, db, secret)
	if err != nil {
		log.Fatal(err)
	}

    // validate with check authorized
	if err := acl.CheckAuthorized(ctx, "goldbank", "*", "list"); err != nil {
		fmt.Println("Not Authorized")
		return
	}
	if err := acl.CheckAuthorized(ctx, "silverbank", "*", "write"); err != nil {
		fmt.Println("Not Authorized")
		return
	}

    // follow the flow of your rules
	fmt.Println("Keep the flow")
```
## TODO

* [x] - Migrate rules from memory to database
* [x] - Organize code style
* [x] - Bind SQL dataset for json call
* [x] - CRUD policies, rules and token
* [x] - Repository for database

## Versioning and license

We use SemVer for versioning. You can see the versions available by checking the tags on this repository.

For more details about our license model, please take a look at the [LICENSE](LICENSE) file

---

2021, thiagozs