package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
			"email":    email,
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
		client := &http.Client{}

		email := req.FormValue("email")
		password := req.FormValue("password")

		// Marshal struct to json
		data, err := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})
		if err != nil {
			log.Fatalf("login error %v\n", err)
		}

		responseBuffer := bytes.NewBuffer(data)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users/login", responseBuffer)
		if err != nil {
			log.Fatalf("error response occured %v\n", err)
		}

		req.Header.Set("Content-Type", "application/json")
		c := req.Cookies()
		fmt.Printf("cookies: %s\n", c)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error repsonse body occurred %v\n", err)
		}

		tokenString := string(respBody)
		tokenString = strings.TrimSuffix(tokenString, "\n")
		tokenString = strings.ReplaceAll(tokenString, "\"", "")

		cookie := &http.Cookie{
			Name:    "jwt-token",
			Value:   tokenString,
			Expires: time.Now().Add(time.Hour * 1),
		}

		fmt.Println(cookie)

		http.SetCookie(res, cookie)
		http.Redirect(res, req, "/", http.StatusOK)
	}
}
