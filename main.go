package main

import (
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/routes"
	"github.com/mofodox/project-live-app/api/utils"
)

func init() {
	// Table creation in Database
	utils.Db.AutoMigrate(&models.Business{}, &models.User{}) 
}

func main() {
	// Starts the server in router.go
	routes.HandleRoutes(":8080")
}
