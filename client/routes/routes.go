package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/client/controllers"
)

func HandleRoutes(addr string) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", controllers.Home)
	router.HandleFunc("/register", controllers.Register).Methods("POST", "GET")
	router.HandleFunc("/login", controllers.Login).Methods("POST", "GET")

	// Business Handlers
	router.HandleFunc("/business/create", controllers.CreateBusiness).Methods("GET")
	router.HandleFunc("/business/create", controllers.ProcessCreateBusiness).Methods("POST")
	router.HandleFunc("/business/update/{id}", controllers.UpdateBusiness).Methods("GET")
	router.HandleFunc("/business/update/{id}", controllers.ProcessUpdateBusiness).Methods("POST")

	// Category Handlers
	router.HandleFunc("/category/create", controllers.CreateCategoryPage).Methods("GET")
	router.HandleFunc("/category/create", controllers.ProcessCategoryForm).Methods("POST")
	router.HandleFunc("/category/update/{id}", controllers.UpdateCategory).Methods("GET")
	router.HandleFunc("/category/update/{id}", controllers.ProcessUpdateCategory).Methods("POST")

	fs := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/css/").Handler(fs)

	log.Printf("Starting server on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
