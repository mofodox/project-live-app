package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/controllers/users"
)

func HandleRoutes(addr string) {
	defaultURI := "/api/v1"

	router := mux.NewRouter().StrictSlash(true)

	/**
	* User routes
	* TODO: Change the REST API Verbs
	*/

	router.HandleFunc(defaultURI+"/users", users.CreateUser).Methods("GET")
	router.HandleFunc(defaultURI+"/users", users.FindAllUsers).Methods("GET")
	router.HandleFunc(defaultURI+"/users/{id}", users.FindUser).Methods("GET")
	router.HandleFunc(defaultURI+"/users/{id}", users.UpdateUser).Methods("GET")
	router.HandleFunc(defaultURI+"/users/{id}", users.DeleteUser).Methods("GET")

	
	log.Printf("Server is running on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
