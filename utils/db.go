package utils

import (
	"gorm.io/gorm"
	"gofp/controllers"
	"gofp/routes"
)
func SetDB(db *gorm.DB) {
	controllers.DB = db
	routes.DB = db
}