package database

import (
	"acl-test-go/model"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseRepo interface {
	LogMode(logger logger.LogLevel)
	CreateTable(models interface{}) error
	CreateTables(models []interface{}) error
	CreateRule(resource, path string, capabilities []string) model.Rule
	CreatePolicy(name string, rules []model.Rule) model.Policie
	CreateToken(privileged bool, userId uint, policies []string) string
	GetPermissonsByToken(token string) (model.TokenResponse, []model.PolicieResponse, error)
}

type databaseRepo struct {
	db     *gorm.DB
	models []interface{}
}

func NewConnection() (DatabaseRepo, error) {
	//db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{
	//Logger: logger.Default.LogMode(logger.Info),
	//DisableForeignKeyConstraintWhenMigrating: true,
	//})
	// just for sqlite3
	//db.Exec("PRAGMA foreign_keys = ON;")

	conn := "root:secret@tcp(127.0.0.1:3306)/openbank?charset=utf8&parseTime=True"
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqldb.SetMaxIdleConns(10)
	sqldb.SetMaxOpenConns(100)
	sqldb.SetConnMaxLifetime(time.Second * 30)

	return &databaseRepo{
		db:     db,
		models: []interface{}{},
	}, nil
}

func (d *databaseRepo) LogMode(logger logger.LogLevel) {
	d.db.Logger.LogMode(logger)
}

func (d *databaseRepo) CreateTable(model interface{}) error {
	return d.db.Migrator().AutoMigrate(&model)
}

func (d *databaseRepo) CreateTables(models []interface{}) error {
	return d.db.Migrator().AutoMigrate(models...)
}

func (d *databaseRepo) CreateRule(resource, path string, capabilities []string) model.Rule {
	pp := fmt.Sprintf(`{"cap":["%s"]}`, strings.Join(capabilities, "\",\""))

	rule := model.Rule{
		Resource:     resource,
		Path:         path,
		Capabilities: datatypes.JSON([]byte(pp)),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	d.db.Create(&rule)
	return rule
}

func (d *databaseRepo) CreatePolicy(name string, rules []model.Rule) model.Policie {
	policy := model.Policie{
		Name:      name,
		Rules:     rules,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	d.db.Create(&policy)
	return policy
}

func (d *databaseRepo) CreateToken(privileged bool, userId uint, policies []string) string {

	pp := fmt.Sprintf(`{"policies":["%s"]}`, strings.Join(policies, "\",\""))
	id := uuid.New()

	token := model.Token{
		UserId:     userId,
		Privileged: privileged,
		Secret:     id.String(),
		Policies:   datatypes.JSON([]byte(pp)),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	d.db.Create(&token)

	return id.String()
}

func (d *databaseRepo) GetPermissonsByToken(token string) (model.TokenResponse, []model.PolicieResponse, error) {
	tokenr := model.Token{}
	policiesr := []model.Policie{}
	rulesr := []model.Rule{}

	// find token by hash
	d.db.Where(&model.Token{Secret: token}).Find(&tokenr)
	tokenResponse, err := bindTokenPolicies(tokenr)
	if err != nil {
		return model.TokenResponse{}, []model.PolicieResponse{}, err
	}

	// get policies by user
	d.db.Where("name IN ?", tokenResponse.Policies).Find(&policiesr)

	// get rules by user
	d.db.Model(model.Rule{}).Joins("INNER JOIN acl_policies ON acl_rules.police_id = acl_policies.id").
		Where("acl_policies.name IN ?", tokenResponse.Policies).Find(&rulesr)
	policies, err := bindPoliciesAndRules(policiesr, rulesr)
	if err != nil {
		return model.TokenResponse{}, []model.PolicieResponse{}, err
	}

	return tokenResponse, policies, nil
}

func getPolicies(hook string, token model.Token) ([]string, error) {
	policies := []string{}
	rawPolicies := map[string]interface{}{}
	if err := json.Unmarshal(token.Policies, &rawPolicies); err != nil {
		return policies, err
	}
	val, ok := rawPolicies[hook]
	if ok {
		for _, v := range val.([]interface{}) {
			policies = append(policies, v.(string))
		}
	}
	return policies, nil
}

func getCapabilities(hook string, rule model.Rule) ([]string, error) {
	caps := []string{}
	rawCap := map[string]interface{}{}
	if err := json.Unmarshal(rule.Capabilities, &rawCap); err != nil {
		return caps, err
	}
	val, ok := rawCap[hook]
	if ok {
		for _, v := range val.([]interface{}) {
			caps = append(caps, v.(string))
		}
	}
	return caps, nil
}

func bindPoliciesAndRules(policiesr []model.Policie, rulesr []model.Rule) ([]model.PolicieResponse, error) {

	policier := []model.PolicieResponse{}
	pop := map[string]*model.PolicieResponse{}

	for _, vv := range policiesr {
		// copy policies struct for policies response
		pol := model.PolicieResponse{}
		_ = copier.Copy(&pol, &vv)
		pop[vv.Name] = &pol

		// check rules
		for _, v := range rulesr {
			if v.PoliceId == vv.Id {
				rule := model.RuleResponse{}
				caps, err := getCapabilities("cap", v)
				if err != nil {
					return []model.PolicieResponse{}, err
				}
				// copy rules struct for rules reponse
				if err := copier.Copy(&rule, &v); err != nil {
					return []model.PolicieResponse{}, err
				}
				rule.Capabilities = caps
				pop[vv.Name].Rules = append(pop[vv.Name].Rules, rule)
			}
		}
	}

	// back to the slice
	for _, v := range pop {
		policier = append(policier, *v)
	}

	return policier, nil
}

func bindTokenPolicies(token model.Token) (model.TokenResponse, error) {

	tokenc := model.TokenResponse{}

	tokenPolicies, err := getPolicies("policies", token)
	if err != nil {
		return model.TokenResponse{}, err
	}
	if err := copier.Copy(&tokenc, &token); err != nil {
		return model.TokenResponse{}, err
	}
	tokenc.Policies = tokenPolicies

	return tokenc, nil
}
