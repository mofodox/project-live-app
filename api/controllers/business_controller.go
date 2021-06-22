package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func (server *Server) AddBusiness(res http.ResponseWriter, req *http.Request) {

}

func (server *Server) UpdateBusiness(res http.ResponseWriter, req *http.Request) {

}

func (server *Server) GetBusiness(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	business_id, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Println("non integer business id")
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}

	var business models.Business

	if err := server.DB.First(&business, business_id).Error; err != nil {
		responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
		return
	}

	responses.JSON(res, http.StatusOK, business)
}

func (server *Server) SearchBusiness(res http.ResponseWriter, req *http.Request) {

}
