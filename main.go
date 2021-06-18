package main

import (
	"github.com/mofodox/project-live-app/api/routes"

	"project-live-app/api/models"
	"project-live-app/api/utils"
)

// db settings
var (
	dbHostname   string
	dbPort       string
	dbUsername   string
	dbPassword   string
	dbDatabase   string
	dbConnection string
)

func init() {
	utils.Db.AutoMigrate(&models.Business{})
}

func main() {
	routes.HandleRoutes(":8080")
}
