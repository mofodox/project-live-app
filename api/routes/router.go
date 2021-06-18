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

	router.HandleFunc(defaultURI+"/users", users.CreateUserHandler).Methods("GET")
	router.HandleFunc(defaultURI+"/users", users.GetAllUsersHandler).Methods("GET")
	router.HandleFunc(defaultURI+"/users/{id}", users.GetUserHandler).Methods("GET")
	router.HandleFunc(defaultURI+"/users/{id}", users.UpdateUserHandler).Methods("GET")
	router.HandleFunc(defaultURI+"/users/{id}", users.DeleteUserHandler).Methods("GET")

	
	log.Printf("Server is running on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
