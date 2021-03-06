package database

import (
	"acl-test-go/model"
	"time"

	"github.com/google/uuid"
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
	CreateSecret(privileged bool, userId uint, policies []string) (string, error)
	GetPermissonsBySecret(secret string) (model.Token, []model.Policie, error)
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

func (d *databaseRepo) CreateSecret(privileged bool, userId uint, policies []string) (string, error) {

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

func (d *databaseRepo) GetPermissonsBySecret(token string) (model.Token, []model.Policie, error) {
	tokenr := model.Token{}
	policiesr := []model.Policie{}

	// find token by hash
	d.db.Model(model.Token{}).Where("secret = ?", token).Find(&tokenr)

	// get policies by user
	tx := d.db.Model(model.Policie{}).
		Where("acl_policies.name IN ?", tokenr.Policies.ToStringArr()).
		Joins("JOIN acl_rules ON acl_rules.police_id = acl_policies.id").
		Preload("Rules").
		Find(&policiesr)
	if tx.Error != nil {
		return model.Token{}, []model.Policie{}, tx.Error
	}

	return tokenr, policiesr, nil
}
