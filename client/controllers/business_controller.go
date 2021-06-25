package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/mofodox/project-live-app/api/models"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").ParseGlob("templates/*"))
}

func CreateBusinessPage(res http.ResponseWriter, req *http.Request) {

	fmt.Println("CREATE BUSINESS PAGE")

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

/*
func UpdateBusinessPage(res http.ResponseWriter, req *http.Request) {

	fmt.Println("UPDATE BUSINESS FORM")

	vars := mux.Vars(req)
	business_id, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	var business *models.Business

	if err := server.DB.First(&business, business_id).Error; err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	fmt.Println(business)

	// Anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		Business   *models.Business
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update Business", nil, business, "", "",
	}

	tpl.ExecuteTemplate(res, "updateBusiness.gohtml", payload)
}
*/

func ProcessBusinessPageForm(res http.ResponseWriter, req *http.Request) {

	fmt.Println("POSTTTT BUSINESS PAGE")

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
		}
	}
}
