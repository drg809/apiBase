package controllers

import (
	"github.com/drg809/apiBase/models"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func Migrate() {
	// DROP
	GormDB.Migrator().DropTable(&models.User{})
	GormDB.Migrator().DropTable(&models.Presale{})

	// CREATE
	GormDB.AutoMigrate(&models.User{})
	GormDB.AutoMigrate(&models.Presale{})

}
