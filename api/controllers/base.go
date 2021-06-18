package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

// db settings
var (
	dbHostname   string
	dbPort       string
	dbUsername   string
	dbPassword   string
	dbDatabase   string
	dbConnection string
)

func (server *Server) Initialize() {
	dbHostname = os.Getenv("MYSQL_HOSTNAME")
	dbPort = os.Getenv("MYSQL_PORT")
	dbUsername = os.Getenv("MYSQL_USERNAME")
	dbPassword = os.Getenv("MYSQL_PASSWORD")
	dbDatabase = os.Getenv("MYSQL_DATABASE")

	dbConnection = dbUsername + ":" + dbPassword + "@tcp(" + dbHostname + ":" + dbPort + ")/" + dbDatabase

	db, err := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	server.DB = db
	server.DB.Debug().AutoMigrate(&models.Business{}, &models.User{}, &models.Category{})
	server.Router = mux.NewRouter().StrictSlash(true)
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Printf("Server is listening on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
