package db

import (
	"errors"
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
