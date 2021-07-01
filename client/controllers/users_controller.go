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

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
	"github.com/mofodox/project-live-app/client/lib"
)

func Register(res http.ResponseWriter, req *http.Request) {

	// Get User
	_, err := lib.IsLoggedIn(req)

	if err == nil {
		// already logged in
		http.Redirect(res, req, "/business", http.StatusSeeOther)
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

		req, err = http.NewRequest(http.MethodPost, lib.ApiBaseURL+"/users", responseBuffer)
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
					Name:     "jwt-token",
					Value:    tokenString,
					Expires:  time.Now().Add(time.Hour * 1),
					HttpOnly: true,
				}

				http.SetCookie(res, cookie)
				http.Redirect(res, req, "/", http.StatusFound)
				return
			}

			http.Redirect(res, req, "/business", http.StatusSeeOther)
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

	lib.Tpl.ExecuteTemplate(res, "register.gohtml", payload)
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

	req, err := http.NewRequest(http.MethodPost, lib.ApiBaseURL+"/users/login", responseBuffer)
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
				Name:     "jwt-token",
				Value:    tokenString,
				Expires:  time.Now().Add(time.Hour * 1),
				HttpOnly: true,
			}

			http.SetCookie(res, cookie)
			http.Redirect(res, req, "/business", http.StatusSeeOther)
			return
		}

		payload.ErrorMsg = err.Error()
	}

	lib.Tpl.ExecuteTemplate(res, "login.gohtml", payload)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:    "jwt-token",
		Value:   "tokenString",
		Expires: time.Now().Add(-time.Hour * 1),
	}

	http.SetCookie(res, cookie)
	http.Redirect(res, req, "/business", http.StatusSeeOther)
}

func ShowProfile(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(res, req, "/business", http.StatusSeeOther)
		return
	}

	user, err := lib.IsLoggedIn(req)
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		ErrorMsg   string
		SuccessMsg string
	}{
		"View Profile", user, "", "",
	}

	client := &http.Client{}

	request, _ := http.NewRequest(http.MethodGet, lib.ApiBaseURL+"/users/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("error sending get user request", err)
		http.Redirect(res, req, "/business", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode == 200 {
		var user *models.User
		marshalErr := json.Unmarshal(data, &user)

		if marshalErr != nil {
			fmt.Println("error unmarshaling at getting user", marshalErr)
			http.Redirect(res, req, "/business", http.StatusTemporaryRedirect)
			return
		}

		payload.User = user

		lib.Tpl.ExecuteTemplate(res, "viewProfile.gohtml", payload)
		fmt.Println("GET View Profile payload", payload.User)
		return
	} else {
		// handle error
		fmt.Println(string(data))
		http.Redirect(res, req, "/business", http.StatusTemporaryRedirect)
		return
	}
}

func UpdateProfile(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	user, err := lib.IsLoggedIn(req)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update Profile", user, "", "",
	}

	client := &http.Client{}

	request, _ := http.NewRequest(http.MethodGet, lib.ApiBaseURL+"/users/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		fmt.Println("error sending get user request", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		var user *models.User
		marshalErr := json.Unmarshal(data, &user)

		if marshalErr != nil {
			fmt.Println("error unmarshaling at update user", marshalErr)
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}

		payload.User = user

		lib.Tpl.ExecuteTemplate(res, "updateUser.gohtml", payload)
		fmt.Println("GET Update user payload", payload.User)
		return
	} else {
		// handle error
		fmt.Println(string(data))
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}
}

func ProcessUpdateProfile(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(res, req, "/users/"+vars["id"], http.StatusSeeOther)
		return
	}

	user, err := lib.IsLoggedIn(req)
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	fullname := req.FormValue("fullname")
	email := req.FormValue("email")
	password := req.FormValue("password")

	user.Fullname = fullname
	user.Email = email
	user.Password = password

	payload := struct {
		PageTitle  string
		User       *models.User
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update User", user, "", "",
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error marshalling at process update user", err)
		payload.ErrorMsg = "An unexpected error has occured while updating user. Please try again."
		lib.Tpl.ExecuteTemplate(res, "updateUser.gohtml", payload)
		return
	}

	client := &http.Client{}

	request, _ := http.NewRequest(http.MethodPut, lib.ApiBaseURL+"/users/"+vars["id"], bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lib.GetJWT(req))

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("error sending process update user request", err)
		payload.ErrorMsg = "An unexpected error has occured while updating user. Please try again."
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		http.Redirect(res, req, "/users/"+vars["id"], http.StatusFound)
		return
	} else {
		respData, _ := ioutil.ReadAll(response.Body)
		var errorResponse responses.ErrorResponse
		json.Unmarshal(respData, &errorResponse)
		payload.ErrorMsg = errorResponse.Error

		if strings.Contains(payload.ErrorMsg, "Duplicate entry") {
			payload.ErrorMsg = "There's already a user with the name provided. Please use a different name."
		}
	}

	lib.Tpl.ExecuteTemplate(res, "updateUser.gohtml", payload)
	fmt.Println("PUT Update user payload", payload.User)
}
