package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func (server *Server) AddComment(res http.ResponseWriter, req *http.Request) {
	// needs to be logged into an account to add comment
	if req.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(req.Body)
		if err == nil {
			newComment := models.Comment{}
			err := json.Unmarshal(reqBody, &newComment)
			if err != nil {
				responses.ERROR(res, http.StatusInternalServerError, err)
				return
			}
			newComment.CreatedAt = time.Now()
			newComment.UpdatedAt = time.Now()
			// need to add User ID and Business ID as well
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

}

func (server *Server) UserComments(res http.ResponseWriter, req *http.Request) {

}

func (server *Server) EditComments(res http.ResponseWriter, req *http.Request) {
	// check whether user is logged in
	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		comment_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}
		var currentComment models.Comment
		err = server.DB.First(&currentComment, comment_ID).Error
		if err != nil {
			responses.ERROR(res, http.StatusNotFound, err)
			return
		}
		if true /*currentComment.UserID == currentUserID */ {
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
				// update inside database
			}
		} else {
			responses.ERROR(res, http.StatusForbidden, errors.New("user not"+
				" authorized to edit this comment"))
		}
	}
}

func (server *Server) RemoveComments(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-type") == "application/json" {
		vars := mux.Vars(req)
		comment_ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}
		var currentComment models.Comment
		err = server.DB.First(&currentComment, comment_ID)

	}
}
