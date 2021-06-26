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
var apiBaseURL string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env values %v\n", err)
	} else {
		log.Println("Successfully loaded the env values")
	}

	apiBaseURL = os.Getenv("APIServerHostname") + ":" + os.Getenv("APIServerPort") + os.Getenv("APIServerBasePath")
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
	request, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/businesses", respBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("Business creation failed")
		http.Redirect(res, req, "/", response.StatusCode)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// Success
		if response.StatusCode == 201 {
			var business *models.Business
			marshalErr := json.Unmarshal(data, &business)

			if marshalErr != nil {
				fmt.Println("Error decoding json at process create business")
				http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
				return
			}

			fmt.Println("Business created successfully")
			fmt.Println(string(data))
			http.Redirect(res, req, "/business/"+strconv.FormatUint(uint64(business.ID), 10), http.StatusSeeOther)
			return
		} else {
			// handle error
			fmt.Println("Business creation failed")
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
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
	request, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/businesses/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("Fetch business at update business form failed")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// success
		if response.StatusCode == 200 {
			marshalErr := json.Unmarshal(data, &business)

			if marshalErr != nil {
				fmt.Println("Error decoding json at update business")
				http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
				return
			}

			payload.Business = business

			tpl.ExecuteTemplate(res, "updateBusiness.gohtml", payload)
			return
		} else {
			// handle error
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}
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
	request, _ := http.NewRequest(http.MethodPut, apiBaseURL+"/businesses/"+vars["id"], respBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("Fetch business at process update business failed")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// success
		if response.StatusCode == 200 {
			http.Redirect(res, req, "/business/"+vars["id"], http.StatusSeeOther)
			return
		} else {
			// handle error
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
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
	request, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/businesses/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("Fetch business at view business failed")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// success
		if response.StatusCode == 200 {
			marshalErr := json.Unmarshal(data, &business)

			if marshalErr != nil {
				fmt.Println("Error decoding json at view business")
				http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
				return
			}

			payload.PageTitle = business.Name
			payload.Business = business

			tpl.ExecuteTemplate(res, "viewBusiness.gohtml", payload)
			return
		} else {
			// handle error
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}
	}
}
