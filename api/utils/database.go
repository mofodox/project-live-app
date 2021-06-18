package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

var Db *gorm.DB

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

	db, err := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	Db = db
}
