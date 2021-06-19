package api

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mofodox/project-live-app/api/controllers"
	"github.com/mofodox/project-live-app/api/seed"
)

var server = controllers.Server{}

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env values %v\n", err)
	} else {
		log.Println("Successfully loaded the env values")
	}

	server.Initialize()
	seed.Load(server.DB)
	server.Run(":8080")
}
