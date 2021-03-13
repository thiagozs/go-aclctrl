package database

import (
	"acl-test-go/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo interface {
	LogMode(logger logger.LogLevel)
	CreateTable(models interface{}) error
	CreateTables(models []interface{}) error
	CreateRule(resource, path string, capabilities []string) model.Rule
	CreatePolicy(name string, userId uint, rules []model.Rule) model.Policie
	CreateSecret(privileged bool, userId uint, policies []string) (string, error)
	GetPermissonsBySecret(secret string) ([]model.Token, []model.Policie, error)
}

type repo struct {
	db     *gorm.DB
	models []interface{}
}

func NewConnection() (Repo, error) {
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

	return &repo{
		db:     db,
		models: []interface{}{},
	}, nil
}

func (d *repo) LogMode(logger logger.LogLevel) {
	d.db.Logger.LogMode(logger)
}

func (d *repo) CreateTable(model interface{}) error {
	return d.db.Migrator().AutoMigrate(&model)
}

func (d *repo) CreateTables(models []interface{}) error {
	return d.db.Migrator().AutoMigrate(models...)
}

func (d *repo) CreateRule(resource, path string, capabilities []string) model.Rule {
	rule := model.Rule{
		Resource:     resource,
		Path:         path,
		Capabilities: capabilities,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	d.db.Create(&rule)
	return rule
}

func (d *repo) CreatePolicy(name string, userId uint, rules []model.Rule) model.Policie {

	policy := model.Policie{
		TokenID:   userId,
		Name:      name,
		Rules:     rules,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	d.db.Create(&policy)
	return policy
}

func (d *repo) CreateSecret(privileged bool, userId uint, policies []string) (string, error) {

	id := uuid.New()

	token := model.Token{
		UserId:     userId,
		Privileged: privileged,
		Secret:     id.String(),
		Policies:   policies,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if tx := d.db.Create(&token); tx.Error != nil {
		return "", tx.Error
	}

	return id.String(), nil
}

func (d *repo) GetPermissonsBySecret(token string) ([]model.Token, []model.Policie, error) {
	secrets := []model.Token{}
	policies := []model.Policie{}

	// find token by hash
	if tx := d.db.Model(model.Token{}).Where("secret IN ?", []string{token}).Find(&secrets); tx.Error != nil {
		return secrets, policies, tx.Error
	}

	// get policies and rules by user
	if tx := d.db.Model(model.Policie{}).
		Where("name IN ?", secrets[0].Policies.ToStringArr()).
		Find(&policies); tx.Error != nil {
		return secrets, policies, tx.Error
	}

	// get rules by policies id
	for idx, v := range policies {
		rules := []model.Rule{}
		if tx := d.db.Model(model.Rule{}).Where("police_id = ?", v.ID).Find(&rules); tx.Error != nil {
			return secrets, policies, tx.Error
		}
		policies[idx].Rules = rules
	}

	return secrets, policies, nil
}
