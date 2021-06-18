package main

import (
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/routes"
	"github.com/mofodox/project-live-app/api/utils"
)

func init() {
	utils.Db.AutoMigrate(&models.Business{}, &models.User{})
}

func main() {
	routes.HandleRoutes(":8080")
}
