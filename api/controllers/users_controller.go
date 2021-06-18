package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func (server *Server) Register(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	
	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	
	responses.JSON(res, http.StatusCreated, userCreated)
}

func (server *Server) Login(res http.ResponseWriter, req *http.Request) {
	// TODO: Login operation
}

func (server *Server) GetUsers(res http.ResponseWriter, req *http.Request) {
	// TODO: Get all users operation
}
