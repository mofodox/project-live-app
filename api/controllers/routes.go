package controllers

func (server *Server) initializeRoutes() {
	defaultURI := "/api/v1"

	/**
	* User routes
	* TODO: Change the REST API Verbs
	 */
	server.Router.HandleFunc(defaultURI+"/users", server.Register).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users/login", server.Login).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users", server.GetUsers).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/users/logout", server.Logout).Methods("POST")

	/**
	* Business routes
	 */
	server.Router.HandleFunc(defaultURI+"/businesses/test", server.TestGeocode).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/businesses", server.SearchBusinesses).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/businesses/{id:[0-9]+}", server.GetBusiness).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/businesses", server.CreateBusiness).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/businesses/{id:[0-9]+}", server.UpdateBusiness).Methods("PUT")

	/**
	* Category routes
	 */
	server.Router.HandleFunc(defaultURI+"/categories/", server.CreateCategory).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/categories/{id:[0-9]+}", server.GetCategory).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/categories/{id:[0-9]+}", server.DeleteCategory).Methods("DELETE")

	/**
	* API Health check route
	 */
	server.Router.HandleFunc(defaultURI+"/health", server.Health).Methods("GET")
}
