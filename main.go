package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mofodox/project-live-app/api/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"project-live-app/api/models"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		dbHostname = os.Getenv("MYSQL_HOSTNAME")
		dbPort = os.Getenv("MYSQL_PORT")
		dbUsername = os.Getenv("MYSQL_USERNAME")
		dbPassword = os.Getenv("MYSQL_PASSWORD")
		dbDatabase = os.Getenv("MYSQL_DATABASE")

		dbConnection = dbUsername + ":" + dbPassword + "@tcp(" + dbHostname + ":" + dbPort + ")/" + dbDatabase
	}

	db, _ := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})

	db.AutoMigrate(&models.Business{})
}

func main() {
	routes.HandleRoutes(":8080")
}
