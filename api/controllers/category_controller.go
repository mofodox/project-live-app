package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func (server *Server) SearchCategory(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		name := req.FormValue("name")

		var categories = []models.Category{}

		result := server.DB.Limit(10)

		if name != "" {
			result = result.Where("name LIKE ?", "%"+name+"%")
		}

		result = result.Find(&categories).Order("name")

		if result.Error != nil {
			responses.ERROR(res, http.StatusInternalServerError, result.Error)
			return
		}

		responses.JSON(res, http.StatusOK, categories)
		return
	}
	responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
}

func (server *Server) GetAllCategory(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {

		var category []models.Category

		if err := server.DB.Find(&category).Error; err != nil {
			responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
			return
		}

		responses.JSON(res, http.StatusOK, category)
		return
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
}

func (server *Server) CreateCategory(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {
		fmt.Println("CREATING CATEGORY WITH CORRECT CONTENT TYPE")
		var userId uint32 = 1

		/*
			// check JWT and get user id
			userId, err := auth.ExtractTokenID(req)
			if err != nil {
				responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}
		*/
		var newCategory models.Category
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newCategory)
			newCategory.UserID = userId

			result := server.DB.Create(&newCategory)

			if result.Error != nil {
				responses.ERROR(res, http.StatusInternalServerError, result.Error)
				return
			}

			responses.JSON(res, http.StatusCreated, newCategory)
			return
		}
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
}

func (server *Server) GetCategory(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		category_id, err := strconv.Atoi(vars["id"])

		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		var category models.Category

		if err := server.DB.First(&category, category_id).Error; err != nil {
			responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
			return
		}

		responses.JSON(res, http.StatusOK, category)
		return
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
}

func (server *Server) DeleteCategory(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		// check JWT and get user id
		userId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		vars := mux.Vars(req)
		category_id, err := strconv.Atoi(vars["id"])

		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		var category models.Category

		if err := server.DB.First(&category, category_id).Error; err != nil {
			responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
			return
		}
		category.UserID = userId

		if err := server.DB.Delete(&category, category_id).Error; err != nil {
			responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
			return
		}

		responses.JSON(res, http.StatusOK, category)
		return
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
}

func (server *Server) UpdateCategory(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		userId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		vars := mux.Vars(req)
		category_id, err := strconv.Atoi(vars["id"])

		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		var category models.Category

		if err := server.DB.First(&category, category_id).Error; err != nil {
			responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
			return
		}
		var updatedCategory models.Category

		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &updatedCategory)

			category.UserID = userId
			category.Name = updatedCategory.Name
			category.Description = updatedCategory.Description

			result := server.DB.Save(&category)

			if result.Error != nil {
				responses.ERROR(res, http.StatusInternalServerError, result.Error)
				return
			}

			responses.JSON(res, http.StatusOK, category)
			return
		}
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("category not found"))
}
