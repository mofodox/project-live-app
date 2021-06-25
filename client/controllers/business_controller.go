package controllers

import (
	"fmt"
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

func ProcessBusinessPageForm(res http.ResponseWriter, req *http.Request) {

	businessName := req.FormValue("name")
	description := req.FormValue("description")
	address := req.FormValue("address")
	zipcode := req.FormValue("zipcode")
	unitno := req.FormValue("unitno")
	website := req.FormValue("website")
	instagram := req.FormValue("instagram")
	facebook := req.FormValue("facebook")

	var newBusiness *models.Business

	newBusiness.Name = businessName
	newBusiness.Description = description
	newBusiness.Address = address
	newBusiness.Zipcode = zipcode
	newBusiness.UnitNo = unitno
	newBusiness.Website = website
	newBusiness.Instagram = instagram
	newBusiness.Facebook = facebook

	newBusiness.Geocode()

	server.DB.Create(&newBusiness)

	// Redirect to Index Page
	http.Redirect(res, req, "/", http.StatusOK)
}
*/
