package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/client/controllers"
)

func HandleRoutes(addr string) {

	router := mux.NewRouter().StrictSlash(true)

	fs := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/css/").Handler(fs)
	router.PathPrefix("/files/").Handler(fs)

	router.HandleFunc("/", controllers.Home)
	router.HandleFunc("/register", controllers.Register).Methods("POST", "GET")
	router.HandleFunc("/login", controllers.Login).Methods("POST", "GET")
	router.HandleFunc("/logout", controllers.Logout).Methods("GET")
	router.HandleFunc("/users/edit/{id}", controllers.UpdateProfile).Methods("GET")
	router.HandleFunc("/users/edit/{id}", controllers.ProcessUpdateProfile).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.ShowProfile).Methods("GET")

	// Business Handlers
	router.HandleFunc("/business", controllers.ListBusiness).Methods("GET")

	router.HandleFunc("/business/{id:[0-9]+}", controllers.ViewBusiness).Methods("GET")
	router.HandleFunc("/business/create", controllers.CreateBusiness).Methods("GET")
	router.HandleFunc("/business/create", controllers.ProcessCreateBusiness).Methods("POST")
	router.HandleFunc("/business/update/{id:[0-9]+}", controllers.UpdateBusiness).Methods("GET")
	router.HandleFunc("/business/update/{id:[0-9]+}", controllers.ProcessUpdateBusiness).Methods("POST")

	// Category Handlers
	router.HandleFunc("/category", controllers.ListCategory).Methods("GET")
	router.HandleFunc("/category/create", controllers.CreateCategoryPage).Methods("GET")
	router.HandleFunc("/category/create", controllers.ProcessCategoryForm).Methods("POST")
	router.HandleFunc("/category/update/{id}", controllers.UpdateCategory).Methods("GET")
	router.HandleFunc("/category/update/{id}", controllers.ProcessUpdateCategory).Methods("POST")
	router.HandleFunc("/category/{id:[0-9]+}", controllers.ViewCategory).Methods("GET")

	// Comment Handlers
	router.HandleFunc("/business/{id}/comment", controllers.CreateComment).Methods("GET")
	router.HandleFunc("/business/{id}/comment", controllers.ProcessCommentForm).Methods("POST")

	log.Printf("Starting server on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
