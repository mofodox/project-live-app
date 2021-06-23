package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Register(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") == "application/json" {
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
		err = user.Validate("")
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}
		
		userCreated, err := user.CreateUser(server.DB)
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}
		
		res.Header().Set("Location", fmt.Sprintf("%s%s/%d", req.Host, req.RequestURI, userCreated.ID))
		responses.JSON(res, http.StatusCreated, userCreated)
	}

	responses.ERROR(res, http.StatusInternalServerError, errors.New("internal server error"))
}

func (server *Server) Login(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") == "application/json" {
		user := &models.User{}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		user.Prepare()
		err = user.Validate("login")
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		token, err := server.SignInUser(res, user.Email, user.Password)
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		responses.JSON(res, http.StatusOK, token)
	}

	responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
}

func (server *Server) SignInUser(res http.ResponseWriter, email, password string) (string, error) {
	user := &models.User{}
	
	var err error

	if err = server.DB.Debug().Model(&models.User{}).Where("email = ?", email).Take(user).Error; err != nil {
		return "", err
	}

	if err = models.VerifyPassword(user.Password, password); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateToken(res, uint32(user.ID))
}

func (server *Server) Logout(res http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie {
		Name: "jwt-token",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(res, cookie)

	responses.JSON(res, http.StatusNoContent, nil)
}

func (server *Server) GetUsers(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") == "application/json" {
		user := models.User{}
	
		users, err := user.FindAllUsers(server.DB)
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		responses.JSON(res, http.StatusOK, users)
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("user not found"))
}

func (server *Server) GetUserById(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") == "application/json" {
		vars := mux.Vars(req)
		uid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			responses.ERROR(res, http.StatusBadRequest, err)
			return
		}

		user := &models.User{}
		userId, err := user.FindUserByID(server.DB, uint32(uid))
		if err != nil {
			responses.ERROR(res, http.StatusBadRequest, err)
			return
		}

		responses.JSON(res, http.StatusOK, userId)
	}

	responses.ERROR(res, http.StatusNotFound, errors.New("user not found"))
}

func (server *Server) UpdateUserById(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") == "application/json" {
		user := models.User{}
	
		vars := mux.Vars(req)
		uid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			responses.ERROR(res, http.StatusBadRequest, err)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		tokenId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		if tokenId != uint32(uid) {
			responses.ERROR(res, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}

		user.Prepare()
		err = user.Validate("update")
		if err != nil {
			responses.ERROR(res, http.StatusUnprocessableEntity, err)
			return
		}

		updatedUser, err := user.UpdateUserByID(server.DB, uint32(uid))
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		responses.JSON(res, http.StatusOK, updatedUser)
	}

	responses.ERROR(res, http.StatusUnprocessableEntity, errors.New("error updating user"))
}

func (server *Server) DeleteUserByID(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") == "applcation/json" {
		vars := mux.Vars(req)
	
		user := models.User{}

		uid, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			responses.ERROR(res, http.StatusBadRequest, err)
			return
		}

		tokenId, err := auth.ExtractTokenID(req)
		if err != nil {
			responses.ERROR(res, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		if tokenId != 0 && tokenId != uint32(uid) {
			responses.ERROR(res, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
			return
		}

		_, err = user.DeleteUser(server.DB, uint32(uid))
		if err != nil {
			responses.ERROR(res, http.StatusInternalServerError, err)
			return
		}

		res.Header().Set("Entity", fmt.Sprintf("%d\n", uid))
		responses.JSON(res, http.StatusNoContent, "success")
	}

	responses.ERROR(res, http.StatusUnprocessableEntity, errors.New("error deleting user"))
}
