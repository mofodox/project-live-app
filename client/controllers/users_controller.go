package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

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
		res.Header().Set("Content-Type", "application/json")
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

		respBody := bytes.NewBuffer(data)

		response, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", respBody)
		if err != nil {
			log.Fatalf("error response occurred %v\n", err)
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("error response body occurred %v\n", err)
		}

		err = json.Unmarshal(body, &models.User{})
		if err != nil {
			log.Printf("error json user{} %v\n", err)
			return
		}

		log.Printf("user registered")

		http.Redirect(res, req, "/", http.StatusCreated)
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		t, err := template.ParseFiles("templates/login.gohtml")
		if err != nil {
			log.Fatalf("template parse file error %v\n", err)
		}

		t.Execute(res, nil)
	}

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
