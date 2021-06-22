package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mofodox/project-live-app/api/models"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize() {
	var err error

	var DbDriver = "mysql"
	var DbHost = os.Getenv("MYSQL_HOSTNAME")
	var DbPort = os.Getenv("MYSQL_PORT")
	var DbUser = os.Getenv("MYSQL_USERNAME")
	var DbPassword = os.Getenv("MYSQL_PASSWORD")
	var DbName = os.Getenv("MYSQL_DATABASE")

	// dbConnection = dbUsername + ":" + dbPassword + "@tcp(" + dbHostname + ":" + dbPort + ")/" + dbDatabase
	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	server.DB, err = gorm.Open(DbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", DbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", DbDriver)
	}

	server.DB.Debug().AutoMigrate(&models.User{})
	server.Router = mux.NewRouter().StrictSlash(true)
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	log.Printf("Server is listening on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
