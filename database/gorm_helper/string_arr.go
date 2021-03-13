package gorm_helper

import (
	"database/sql/driver"
	"errors"
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
	switch value := src.(type) {
	case string:
		*s = strings.Split(value, "|")
		return nil
	case []byte:
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

func (s *StringArr) ToStringArr() []string {
	return *s
}
