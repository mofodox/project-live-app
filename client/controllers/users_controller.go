package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mofodox/project-live-app/api/models"
)

func Register(res http.ResponseWriter, req *http.Request) {
	// Anonymous payload
	payload := struct {
		PageTitle  string
		ErrorMsg   string
		SuccessMsg string
	}{
		"User Register", "", "",
	}

	tpl.ExecuteTemplate(res, "register.gohtml", payload)

	if req.Method == http.MethodPost {
		client := &http.Client{}

		fullname := req.FormValue("fullname")
		email := req.FormValue("email")
		password := req.FormValue("password")

		data, err := json.Marshal(map[string]string{
			"fullname": fullname,
			"email": email,
			"password": password,
		})
		if err != nil {
			log.Fatalf("register error %v\n", err)
		}

		responseBuffer := bytes.NewBuffer(data)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", responseBuffer)
		if err != nil {
			log.Fatalf("error response occurred %v\n", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(respBody, &models.User{}); err != nil {
			log.Fatal(err)
		}
		log.Println("User registered")

		http.Redirect(res, req, "/", http.StatusCreated)
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	// Anonymous payload
	payload := struct {
		PageTitle  string
		ErrorMsg   string
		SuccessMsg string
	}{
		"User Login", "", "",
	}

	tpl.ExecuteTemplate(res, "login.gohtml", payload)

	if req.Method == http.MethodPost {
		email := req.FormValue("email")
		password := req.FormValue("password")

		data, err := json.Marshal(map[string]string{
			"email": email,
			"password": password,
		})
		if err != nil {
			log.Fatalf("login error %v\n", err)
		}

		respBody := bytes.NewBuffer(data)

		response, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users/login", respBody)
		if err != nil {
			log.Fatalf("error response occured %v\n", err)
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("error repsonse body occurred %v\n", err)
		}

		err = json.Unmarshal(body, &models.User{})
		if err != nil {
			log.Printf("error json user{} %v\n", err)
			return
		}

		log.Printf("user logged in")

		http.Redirect(res, req, "/", http.StatusCreated)
	}
}
