package gorm_helper

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Usage of helper []string
// Just put the tag of gorm for text like that
// -------------------------------------------
// Category StringArr `gorm:"type:text"`
// The gorm lib cannot handle []string, you can
// write a helper like below to handle that.

type StringArr []string

func (s *StringArr) Scan(src interface{}) error {
	fmt.Println("TypeOf = ", reflect.TypeOf(src))
	fmt.Println("ValueOf = ", reflect.ValueOf(src))
	switch value := src.(type) {
	case string:
		fmt.Println("SOU STRING")
		*s = strings.Split(value, "|")
		return nil
	case []byte:
		fmt.Println("SOU []BYTE")
		*s = strings.Split(string(value), "|")
		return nil
	}
	return errors.New("failed to scan stringarr field - source is not a string or bytes")
}

func (s StringArr) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, "|"), nil
}

func (s StringArr) GormDataType() string {
	return "text"
}

func (s StringArr) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "text"
	}
	return ""
}

func (s StringArr) ToStringArr() []string {
	return s
}
