package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/client/controllers"
)

func HandleRoutes(addr string) {
	staticDir := "/public"

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", controllers.Home)
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST", "GET")

	// Business Handlers
	router.HandleFunc("/business/create", controllers.CreateBusinessPage).Methods("GET")

	/*
		router.HandleFunc("/business/create", controllers.CreateBusinessPage).Methods("GET")
		router.HandleFunc("/business/create", controllers.ProcessBusinessPageForm).Methods("POST")
		router.HandleFunc("/business/update/{id}", controllers.UpdateBusinessPage).Methods("GET")
	*/

	router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	log.Printf("Starting server on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
