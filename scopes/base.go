package scopes

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(db *gorm.DB) error {
	if DB != nil {
		return errors.New("db already set")
	}
	DB = db
	return nil
}

type ScopeFunc = func(db *gorm.DB) *gorm.DB