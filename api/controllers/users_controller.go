package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
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
	
	responses.JSON(res, http.StatusCreated, userCreated)
}

func (server *Server) Login(res http.ResponseWriter, req *http.Request) {
	// TODO: Login operation
	res.Header().Set("Content-Type", "application/json")
}

func (server *Server) GetUsers(res http.ResponseWriter, req *http.Request) {
	user := models.User{}
	
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(res, http.StatusOK, users)
}
