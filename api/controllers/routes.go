package controllers

import "github.com/mofodox/project-live-app/api/middlewares"

func (server *Server) initializeRoutes() {
	defaultURI := "/api/v1"

	/**
	* User routes
	 */
	server.Router.HandleFunc(defaultURI+"/users", middlewares.SetMiddlewareJSON(server.Register)).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users/login", middlewares.SetMiddlewareJSON(server.Login)).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users", middlewares.SetMiddlewareJSON(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/users/logout", middlewares.SetMiddlewareJSON(server.Logout)).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users/{id}", middlewares.SetMiddlewareJSON(server.GetUserById)).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(server.UpdateUserById))).Methods("PUT")

	/**
	* Business routes
	 */
	server.Router.HandleFunc(defaultURI+"/businesses", server.SearchBusinesses).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/businesses/{id:[0-9]+}", server.GetBusiness).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/businesses", server.CreateBusiness).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/businesses/{id:[0-9]+}", server.UpdateBusiness).Methods("PUT")
	server.Router.HandleFunc(defaultURI+"/businesses/{id:[0-9]+}", server.DeleteBusiness).Methods("DELETE")

	/**
	* Category routes
	 */
	server.Router.HandleFunc(defaultURI+"/categories/", server.CreateCategory).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/categories/{id:[0-9]+}", server.GetCategory).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/categories/{id:[0-9]+}", server.DeleteCategory).Methods("DELETE")
	server.Router.HandleFunc(defaultURI+"/categories/{id:[0-9]+}", server.UpdateCategory).Methods("PUT")

	/**
	* API Health check route
	 */
	server.Router.HandleFunc(defaultURI+"/health", server.Health).Methods("GET")
}
