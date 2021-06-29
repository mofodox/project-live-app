package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

const BusinessSearchLimit = 5

func (server *Server) CreateBusiness(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {

		// check JWT and get user id
		userId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		var newBusiness *models.Business
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newBusiness)

			// todo: sanitization

			// validation
			err := newBusiness.Validate()

			if err != nil {
				responses.ERROR(res, http.StatusInternalServerError, err)
				return
			}

			newBusiness.UserID = userId

			// get lat / lng automatically
			if newBusiness.Lat == 0 && newBusiness.Lng == 0 {
				// todo: change to go routine
				newBusiness.Geocode()
			}

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

// Setting of status to inactive
func (server *Server) DeleteBusiness(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {

		// check JWT and get user id
		userId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

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

		business.UserID = userId
		business.Status = "inactive"

		result := server.DB.Save(&business)

		if result.Error != nil {
			responses.ERROR(res, http.StatusInternalServerError, result.Error)
			return
		}

		responses.JSON(res, http.StatusOK, business)
		return
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
}

func (server *Server) UpdateBusiness(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {

		// check JWT and get user id
		userId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

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

		var updatedBusiness models.Business

		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &updatedBusiness)

			// todo: sanitization

			// validation
			err := updatedBusiness.Validate()

			if err != nil {
				responses.ERROR(res, http.StatusInternalServerError, err)
				return
			}

			business.UserID = userId
			business.Name = updatedBusiness.Name
			business.Description = updatedBusiness.Description
			business.ShortDescription = updatedBusiness.ShortDescription
			business.Address = updatedBusiness.Address
			business.UnitNo = updatedBusiness.UnitNo
			business.Zipcode = updatedBusiness.Zipcode
			business.Lat = updatedBusiness.Lat
			business.Lng = updatedBusiness.Lng
			business.Status = updatedBusiness.Status

			business.Website = updatedBusiness.Website
			business.Instagram = updatedBusiness.Instagram
			business.Facebook = updatedBusiness.Facebook

			business.Geocode()

			result := server.DB.Save(&business)

			if result.Error != nil {
				responses.ERROR(res, http.StatusInternalServerError, result.Error)
				return
			}

			responses.JSON(res, http.StatusOK, business)
			return
		}
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
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

// http://localhost:8080/api/v1/businesses/search?page=1&q=wayne
func (server *Server) SearchBusinesses(res http.ResponseWriter, req *http.Request) {

	if req.Header.Get("Content-type") == "application/json" {
		q := req.FormValue("q")
		status := strings.ToLower(req.FormValue("status"))
		page := req.FormValue("page")

		// Status
		if status != "" && status != "active" {
			status = "inactive"
		}

		// Pagination
		pageNo, err := strconv.Atoi(page)

		if err != nil || pageNo <= 0 {
			pageNo = 1
		}

		limit := BusinessSearchLimit
		offset := (pageNo - 1) * limit

		var businesses = []*models.Business{}

		// Construct query
		result := server.DB.Offset(offset).Limit(limit)
		countResult := server.DB.Table("businesses")

		// Just using sub string search for now inside business name / description / short description
		if q != "" {
			result = result.Where("name LIKE ?", "%"+q+"%").Or("description LIKE ?", "%"+q+"%").Or("short_description LIKE ?", "%"+q+"%")
			countResult = countResult.Where("name LIKE ?", "%"+q+"%").Or("description LIKE ?", "%"+q+"%").Or("short_description LIKE ?", "%"+q+"%")
		}

		if status != "" {
			result = result.Where("status = ?", status)
			countResult = countResult.Where("status = ?", status)
		}

		var count int
		result = result.Find(&businesses).Order("name")
		_ = countResult.Count(&count)

		if result.Error != nil {
			responses.ERROR(res, http.StatusInternalServerError, result.Error)
			return
		}

		var businessSearchResult models.BusinessSearchResult

		businessSearchResult.Businesses = businesses
		businessSearchResult.Total = count
		businessSearchResult.Limit = BusinessSearchLimit

		responses.JSON(res, http.StatusOK, businessSearchResult)
		return
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("business not found"))
}
