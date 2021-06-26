package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mofodox/project-live-app/api/models"
)

var tpl *template.Template

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env values %v\n", err)
	} else {
		log.Println("Successfully loaded the env values")
	}

	tpl = template.Must(template.New("").ParseGlob("templates/*"))
}

func CreateBusiness(res http.ResponseWriter, req *http.Request) {

	// Anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		ErrorMsg   string
		SuccessMsg string
	}{
		"Create Business", nil, "", "",
	}

	tpl.ExecuteTemplate(res, "createBusiness.gohtml", payload)
}

func ProcessCreateBusiness(res http.ResponseWriter, req *http.Request) {

	businessName := req.FormValue("name")
	description := req.FormValue("description")
	address := req.FormValue("address")
	zipcode := req.FormValue("zipcode")
	unitno := req.FormValue("unitno")
	website := req.FormValue("website")
	instagram := req.FormValue("instagram")
	facebook := req.FormValue("facebook")

	data, err := json.Marshal(map[string]string{
		"name":        businessName,
		"description": description,
		"address":     address,
		"zipcode":     zipcode,
		"unitno":      unitno,
		"website":     website,
		"instagram":   instagram,
		"facebook":    facebook,
	})
	if err != nil {
		log.Fatalf("login error %v\n", err)
	}

	respBody := bytes.NewBuffer(data)

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/businesses", respBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		// handle error
		fmt.Println("Business creation failed")
		http.Redirect(res, req, "/", response.StatusCode)
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)

		// Success
		if response.StatusCode == 201 {
			fmt.Println("Business created successfully")
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusOK)
		} else {
			// handle error
			fmt.Println("Business creation failed")
			http.Redirect(res, req, "/", response.StatusCode)
		}
	}
}

func UpdateBusiness(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	var business *models.Business

	// Anonymous payload
	payload := struct {
		PageTitle  string
		Business   *models.Business
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update Business", nil, "", "",
	}

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/businesses/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		// handle error
		fmt.Println("Business update failed")
		http.Redirect(res, req, "/", response.StatusCode)
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println(response.StatusCode)

		data, _ := ioutil.ReadAll(response.Body)
		marshalErr := json.Unmarshal(data, &business)

		if marshalErr != nil {
			fmt.Println("Error decoding json")
		}

		payload.Business = business

		tpl.ExecuteTemplate(res, "updateBusiness.gohtml", payload)
		return
	}
}

func ProcessUpdateBusiness(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	businessName := req.FormValue("name")
	description := req.FormValue("description")
	address := req.FormValue("address")
	zipcode := req.FormValue("zipcode")
	unitno := req.FormValue("unitno")
	website := req.FormValue("website")
	instagram := req.FormValue("instagram")
	facebook := req.FormValue("facebook")

	data, err := json.Marshal(map[string]string{
		"name":        businessName,
		"description": description,
		"address":     address,
		"zipcode":     zipcode,
		"unitno":      unitno,
		"website":     website,
		"instagram":   instagram,
		"facebook":    facebook,
	})
	if err != nil {
		log.Fatalf("login error %v\n", err)
	}

	respBody := bytes.NewBuffer(data)

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/businesses/"+vars["id"], respBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		// handle error
		fmt.Println("Business creation failed")
		http.Redirect(res, req, "/", response.StatusCode)
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)

		// Success
		if response.StatusCode == 200 {
			fmt.Println("Business updated successfully")
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusOK)
		} else {
			// handle error
			fmt.Println("Business update failed")
			http.Redirect(res, req, "/", response.StatusCode)
		}
	}
}

func ViewBusiness(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	var business *models.Business

	// Anonymous payload
	payload := struct {
		PageTitle   string
		Business    *models.Business
		ErrorMsg    string
		SuccessMsg  string
		GMapsAPIKey string
	}{
		"View Business", nil, "", "", os.Getenv("GMapsPublicAPI"),
	}

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/businesses/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		// handle error
		fmt.Println("Business update failed")
		http.Redirect(res, req, "/", response.StatusCode)
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println(response.StatusCode)

		data, _ := ioutil.ReadAll(response.Body)
		marshalErr := json.Unmarshal(data, &business)

		if marshalErr != nil {
			fmt.Println("Error decoding json")
		}

		payload.PageTitle = business.Name
		payload.Business = business

		tpl.ExecuteTemplate(res, "viewBusiness.gohtml", payload)
		return
	}
}
