package controllers

import (
	"fmt"
	"net/http"
)

func (server *Server) Health(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	fmt.Fprintf(res, "API is OK!")
}
