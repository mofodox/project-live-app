package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func (server *Server) CreateBusiness(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {
		var newBusiness models.Business
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newBusiness)
			result := server.DB.Create(&newBusiness)

			if result.Error != nil {
				responses.ERROR(res, http.StatusInternalServerError, result.Error)
				return
			}

			responses.JSON(res, http.StatusCreated, newBusiness)
			return
		}
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
}

func (server *Server) UpdateBusiness(res http.ResponseWriter, req *http.Request) {

}

func (server *Server) GetBusiness(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		business_id, err := strconv.Atoi(vars["id"])

		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		var business models.Business

		if err := server.DB.First(&business, business_id).Error; err != nil {
			responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
			return
		}

		responses.JSON(res, http.StatusOK, business)
		return
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
}

func (server *Server) SearchBusiness(res http.ResponseWriter, req *http.Request) {

}
