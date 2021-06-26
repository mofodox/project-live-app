package main

import (
	"os"

	"github.com/mofodox/project-live-app/client/routes"
)

func main() {
	routes.HandleRoutes(":" + os.Getenv("ClientServerPort"))
}
