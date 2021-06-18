package api

import (
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}

	/*
	 
	*/
	server.DB.Debug().AutoMigrate(&models.User{})
}
