package rest

import (
	"github.com/sunshine-dev-code/rest/db"
	"gorm.io/gorm"
)

func InitDB(DB *gorm.DB) error {
	return db.InitDB(DB)
}
