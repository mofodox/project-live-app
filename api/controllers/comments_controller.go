package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func (server *Server) AddComment(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		userID, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("no authorization"))
		}
		reqBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			newComment := models.Comment{}
			err := json.Unmarshal(reqBody, &newComment)
			if err != nil {
				responses.ERROR(res, http.StatusInternalServerError, err)
				return
			}
			if newComment.BusinessID == 0 || (newComment.Content == "" /* || reivew == 0 */) {
				responses.ERROR(res, http.StatusBadRequest, errors.New("not enough"+
					" information provided"))
			}
			newComment.UserID = userID
			newComment.CreatedAt = time.Now()
			newComment.UpdatedAt = time.Now()
			err = server.DB.Debug().Create(&newComment).Error
			if err != nil {
				responses.ERROR(res, http.StatusInternalServerError, err)
				return
			}
			responses.JSON(res, http.StatusCreated, newComment)
			return
		}
	}
}

func (server *Server) BusinessComments(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		business_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}
		// Might chance to linked list later
		bComments := []models.Comment{}
		err = server.DB.Where("business_id = ?", business_ID).Find(&bComments).Error
		if err != nil {
			responses.ERROR(res, http.StatusNotFound, err)
			return
		}
		responses.JSON(res, http.StatusOK, bComments)
	}
}

func (server *Server) UserComments(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		user_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
		}
		// Might chance to linked list later
		uComments := []models.Comment{}
		err = server.DB.Where("user_id = ?", user_ID).Find(&uComments).Error
		if err != nil {
			responses.ERROR(res, http.StatusNotFound, err)
			return
		}
		responses.JSON(res, http.StatusOK, uComments)
	}
}

func (server *Server) GetComment(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		comment_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}
		comment := models.Comment{}
		err = server.DB.First(&comment, comment_ID).Error
		if err != nil {
			responses.ERROR(res, http.StatusNotFound, err)
			return
		}
		responses.JSON(res, http.StatusOK, comment)
	}
}

func (server *Server) EditComments(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		userID, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("no authorization"))
			return
		}
		vars := mux.Vars(req)
		comment_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}
		currentComment := models.Comment{}
		err = server.DB.First(&currentComment, comment_ID).Error
		if err != nil {
			responses.ERROR(res, http.StatusNotFound, err)
			return
		}
		if currentComment.UserID == userID {
			updateComment := models.Comment{}
			reqBody, err := ioutil.ReadAll(req.Body)
			if err == nil {
				json.Unmarshal(reqBody, &updateComment)
				if updateComment.Content == currentComment.Content {
					responses.JSON(res, http.StatusNotModified, currentComment)
					return
				}
				currentComment.Content = updateComment.Content
				currentComment.UpdatedAt = time.Now()
				err = server.DB.Save(&currentComment).Error
				if err != nil {
					responses.ERROR(res, http.StatusInternalServerError, err)
					return
				}
				responses.JSON(res, http.StatusOK, currentComment)
			}
		} else {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("no authorization"))
		}
	}
}

func (server *Server) RemoveComments(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		userID, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("no authorization"))
		}
		vars := mux.Vars(req)
		comment_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}
		currentComment := models.Comment{}
		err = server.DB.First(&currentComment, comment_ID).Error
		if err != nil {
			responses.ERROR(res, http.StatusNotFound, err)
		}
		if currentComment.UserID == userID {
			err = server.DB.Delete(&currentComment, comment_ID).Error
			if err != nil {
				responses.ERROR(res, http.StatusInternalServerError, err)
			}
			responses.JSON(res, http.StatusOK, "Deleted successfully")
		} else {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("no authorization"))
		}
	}
}
