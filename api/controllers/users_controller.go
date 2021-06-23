package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Register(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

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
	
	userCreated, err := user.CreateUser(server.DB)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}
	
	responses.JSON(res, http.StatusCreated, userCreated)
}

func (server *Server) Login(res http.ResponseWriter, req *http.Request) {
	// TODO: Login operation
	res.Header().Set("Content-Type", "application/json")

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

	token, err := server.SignInUser(res, user.Email, user.Password)
	if err != nil {
		responses.ERROR(res, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(res, http.StatusOK, token)
}

func (server *Server) SignInUser(res http.ResponseWriter, email, password string) (string, error) {
	user := &models.User{}

	if err := server.DB.Debug().Model(&models.User{}).Where("email = ?", email).Take(user).Error; err != nil {
		return "", err
	}

	if err := models.VerifyPassword(user.Password, password); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Expires in 1 hour
	})

	token, err := claims.SignedString([]byte(os.Getenv("HBB_SECRET_KEY")))
	if err != nil {
		log.Fatalf("token error %s\n", err)
	}

	cookie := &http.Cookie {
		Name: "jwt-token",
		Value: token,
		Expires: time.Now().Add(time.Hour * 1),
		HttpOnly: true,
	}

	http.SetCookie(res, cookie)

	return token, nil
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
	res.Header().Set("Content-Type", "application/json")

	user := models.User{}
	
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(res, http.StatusOK, users)
}

func (server *Server) GetUserById(res http.ResponseWriter, req *http.Request) {
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

func (server *Server) UpdateUserById(res http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	
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

	updatedUser, err := user.UpdateUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(res, http.StatusOK, updatedUser)
}
