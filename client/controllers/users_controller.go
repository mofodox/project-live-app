package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mofodox/project-live-app/api/auth"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func IsLoggedIn(req *http.Request) (*models.User, error) {
	myCookie, err := req.Cookie("jwt-token")
	if err != nil {
		return nil, err
	}

	token := myCookie.Value

	userID, err := auth.GetTokenID(token)
	if err != nil {
		return nil, err
	}

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/users/"+strconv.FormatUint(uint64(userID), 10), nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending get user info request", err)
		return nil, err
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		var user *models.User
		marshalErr := json.Unmarshal(data, &user)

		if marshalErr != nil {
			return nil, marshalErr
		}

		return user, nil
	}

	return nil, errors.New("unable to fetch user")
}

func GetJWT(req *http.Request) string {

	myCookie, err := req.Cookie("jwt-token")
	if err == nil {
		return myCookie.Value
	}

	return ""
}

func Register(res http.ResponseWriter, req *http.Request) {

	// Get User
	_, err := IsLoggedIn(req)

	if err == nil {
		// already logged in
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}

	// Anonymous payload
	payload := struct {
		PageTitle  string
		ErrorMsg   string
		SuccessMsg string
		User       *models.User
	}{
		"User Register", "", "", nil,
	}

	if req.Method == http.MethodPost {

		fullname := req.FormValue("fullname")
		email := req.FormValue("email")
		password := req.FormValue("password")

		data, err := json.Marshal(map[string]string{
			"fullname": fullname,
			"email":    email,
			"password": password,
		})
		if err != nil {
			log.Fatalf("register error %v\n", err)
		}

		responseBuffer := bytes.NewBuffer(data)

		req, err = http.NewRequest(http.MethodPost, apiBaseURL+"/users", responseBuffer)
		if err != nil {
			log.Fatalf("error response occurred %v\n", err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var user models.User

		if err := json.Unmarshal(respBody, &user); err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode == 201 {

			// auto login user
			tokenString, err := sendLoginRequest(email, password)

			if err == nil {
				cookie := &http.Cookie{
					Name:    "jwt-token",
					Value:   tokenString,
					Expires: time.Now().Add(time.Hour * 1),
				}

				http.SetCookie(res, cookie)
				http.Redirect(res, req, "/", http.StatusFound)
				return
			}

			http.Redirect(res, req, "/login", http.StatusSeeOther)
			return
		} else {
			var errorResponse responses.ErrorResponse
			json.Unmarshal(respBody, &errorResponse)
			payload.ErrorMsg = errorResponse.Error

			if strings.Contains(payload.ErrorMsg, "Duplicate") {
				payload.ErrorMsg = "The email address provided is already registered."
			}
		}
	}

	tpl.ExecuteTemplate(res, "register.gohtml", payload)
}

func sendLoginRequest(email string, password string) (string, error) {
	client := &http.Client{}

	// Marshal struct to json
	data, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return "", errors.New("error marshaling login request")
	}

	responseBuffer := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPost, apiBaseURL+"/users/login", responseBuffer)
	if err != nil {
		log.Fatalf("error response occured %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error response body occurred %v\n", err)
	}

	if resp.StatusCode == 200 {
		tokenString := string(respBody)
		tokenString = strings.TrimSuffix(tokenString, "\n")
		tokenString = strings.ReplaceAll(tokenString, "\"", "")

		return tokenString, nil
	}

	return "", errors.New("your account and/or password is incorrect, please try again")
}

func Login(res http.ResponseWriter, req *http.Request) {

	// Anonymous payload
	payload := struct {
		PageTitle  string
		ErrorMsg   string
		SuccessMsg string
		User       *models.User
	}{
		"User Login", "", "", nil,
	}

	if req.Method == http.MethodPost {

		email := req.FormValue("email")
		password := req.FormValue("password")
		tokenString, err := sendLoginRequest(email, password)

		if err == nil {
			cookie := &http.Cookie{
				Name:    "jwt-token",
				Value:   tokenString,
				Expires: time.Now().Add(time.Hour * 1),
			}

			http.SetCookie(res, cookie)
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		payload.ErrorMsg = err.Error()
	}

	tpl.ExecuteTemplate(res, "login.gohtml", payload)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:    "jwt-token",
		Value:   "tokenString",
		Expires: time.Now().Add(-time.Hour * 1),
	}

	http.SetCookie(res, cookie)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}
