package controllers

func (server *Server) initializeRoutes() {
	defaultURI := "/api/v1"

	/**
	* User routes
	* TODO: Change the REST API Verbs
	 */
	server.Router.HandleFunc(defaultURI+"/users", server.Register).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users", server.Login).Methods("POST")
	server.Router.HandleFunc(defaultURI+"/users", server.GetUsers).Methods("GET")

	/**
	* Business routes
	 */
	server.Router.HandleFunc(defaultURI+"/businesses/{id}", server.GetBusiness).Methods("GET")
	server.Router.HandleFunc(defaultURI+"/businesses", server.CreateBusiness).Methods("POST")

	/**
	* API Health check route
	 */
	server.Router.HandleFunc(defaultURI+"/health", server.Health).Methods("GET")
}
