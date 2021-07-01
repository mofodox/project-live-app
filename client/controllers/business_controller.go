package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
	"github.com/mofodox/project-live-app/client/lib"
)

var tpl *template.Template
var apiBaseURL string

const BusinessSearchLimit = 5

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env values %v\n", err)
	} else {
		log.Println("Successfully loaded the env values")
	}

	apiBaseURL = os.Getenv("APIServerHostname") + ":" + os.Getenv("APIServerPort") + os.Getenv("APIServerBasePath")

	funcMap := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},

		"ucFirst": func(str string) string {
			if len(str) == 0 {
				return ""
			}
			tmp := []rune(str)
			tmp[0] = unicode.ToUpper(tmp[0])
			return string(tmp)
		},

		"formatDistance": func(distance float64) string {
			if distance > 1 {
				return fmt.Sprint(math.Floor(distance*100)/100) + "KM"
			} else {
				return fmt.Sprint(math.Floor(distance*100)) + "M"
			}
		},

		"formatCommentDate": func(t time.Time) string {
			return t.Format("2 Jan 2006 3:04 PM")
		},
	}

	tpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*"))
}

func ListBusiness(res http.ResponseWriter, req *http.Request) {

	pageNo := 1
	searchedLocation := ""

	q := req.FormValue("q")
	status := strings.ToLower(req.FormValue("status"))
	page := req.FormValue("page")

	// pagination
	p, err := strconv.Atoi(page)

	if err != nil || p <= 0 {
		pageNo = 1
	} else {
		pageNo = p
	}

	querystring := "?page=" + strconv.Itoa(pageNo)
	addtionalQuerystring := ""

	// status
	if status != "" && status != "active" {
		addtionalQuerystring += "&status=inactive"
	}

	// query
	// hackish way to incorporate location search for now...
	if q != "" {
		if strings.Contains(strings.ToLower(q), "location:") {
			q = strings.ReplaceAll(strings.ToLower(q), "location:", "")
			addtionalQuerystring += "&location=" + q
			searchedLocation = q
		} else {
			addtionalQuerystring += "&q=" + q
		}
	}

	client := &http.Client{}
	url := apiBaseURL + "/businesses" + querystring + addtionalQuerystring
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("Fetch businesses at list business failed")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		//var businesses = []*models.Business{}
		var businessSearchResult models.BusinessSearchResult

		marshalErr := json.Unmarshal(data, &businessSearchResult)

		if marshalErr != nil {
			fmt.Println("Error decoding json at list businesses")
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}

		// anonymous payload
		payload := struct {
			PageTitle  string
			Location   string
			StartNo    int
			PrevURL    string
			NextURL    string
			Total      int
			Businesses []*models.Business
			User       *models.User
			ErrorMsg   string
			SuccessMsg string
		}{
			"Businesses", searchedLocation, 1, "", "", businessSearchResult.Total, businessSearchResult.Businesses, nil, "", "",
		}

		// page limit
		if pageNo > 1 {
			payload.StartNo = (pageNo-1)*businessSearchResult.Limit + 1
		}

		// Pagination

		// Previous Page
		if payload.Total != 0 && pageNo > 1 {
			querystring := "?page=" + strconv.Itoa(pageNo-1) + addtionalQuerystring
			payload.PrevURL = querystring
		}

		// Next Page
		if pageNo*businessSearchResult.Limit < payload.Total {
			querystring := "?page=" + strconv.Itoa(pageNo+1) + addtionalQuerystring
			payload.NextURL = querystring
		}

		// Get User
		user, err := lib.IsLoggedIn(req)

		if err == nil {
			payload.User = user
		}

		tpl.ExecuteTemplate(res, "businessListing.gohtml", payload)
		return
	} else {
		// handle error
		fmt.Println("Business listing failed")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}
}

func CreateBusiness(res http.ResponseWriter, req *http.Request) {

	// Get User
	user, err := lib.IsLoggedIn(req)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		Business   *models.Business
		ErrorMsg   string
		SuccessMsg string
	}{
		"Create Business", user, nil, "", "",
	}

	tpl.ExecuteTemplate(res, "createBusiness.gohtml", payload)
}

func ProcessCreateBusiness(res http.ResponseWriter, req *http.Request) {

	// Get User
	user, err := lib.IsLoggedIn(req)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		Business   *models.Business
		ErrorMsg   string
		SuccessMsg string
	}{
		"Create Business", user, nil, "", "",
	}

	var business models.Business

	business.Name = req.FormValue("name")
	business.ShortDescription = req.FormValue("shortDescription")
	business.Description = req.FormValue("description")
	business.Address = req.FormValue("address")
	business.Zipcode = req.FormValue("zipcode")
	business.UnitNo = req.FormValue("unitno")
	business.Website = req.FormValue("website")
	business.Instagram = req.FormValue("instagram")
	business.Facebook = req.FormValue("facebook")

	payload.Business = &business

	data, err := json.Marshal(business)
	if err != nil {
		fmt.Println("error marshalling at process create business", err)
		payload.ErrorMsg = "An unexpected error has occured while creating business. Please try again."
		tpl.ExecuteTemplate(res, "createBusiness.gohtml", payload)
		return
	}

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, apiBaseURL+"/businesses", bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lib.GetJWT(req))
	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending process create business request")
		payload.ErrorMsg = "An unexpected error has occured while creating business. Please try again."
		tpl.ExecuteTemplate(res, "createBusiness.gohtml", payload)
		return
	}
	defer response.Body.Close()

	respData, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 201 {
		marshalErr := json.Unmarshal(respData, &business)

		if marshalErr != nil {
			fmt.Println("error unmarshaling at process create business", marshalErr)
			payload.ErrorMsg = "An unexpected error has occured while creating business. Please try again."
			tpl.ExecuteTemplate(res, "createBusiness.gohtml", payload)
			return
		}

		http.Redirect(res, req, "/business/"+strconv.FormatUint(uint64(business.ID), 10), http.StatusFound)
		return
	} else {
		// handle error
		var errorResponse responses.ErrorResponse
		json.Unmarshal(respData, &errorResponse)
		payload.ErrorMsg = errorResponse.Error

		if strings.Contains(payload.ErrorMsg, "Duplicate entry") {
			payload.ErrorMsg = "There's already a business with the name provided. Please use a different name."
		}
	}

	tpl.ExecuteTemplate(res, "createBusiness.gohtml", payload)
}

func UpdateBusiness(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Get User
	user, err := lib.IsLoggedIn(req)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		Business   *models.Business
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update Business", user, nil, "", "",
	}

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/businesses/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending get business request", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		var business *models.Business
		marshalErr := json.Unmarshal(data, &business)

		if marshalErr != nil {
			fmt.Println("error unmarshaling at update business", marshalErr)
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

func ProcessUpdateBusiness(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	businessID, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// Get User
	user, err := lib.IsLoggedIn(req)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	}

	var business models.Business

	business.ID = uint32(businessID)
	business.Name = req.FormValue("name")
	business.ShortDescription = req.FormValue("shortDescription")
	business.Description = req.FormValue("description")
	business.Address = req.FormValue("address")
	business.Zipcode = req.FormValue("zipcode")
	business.UnitNo = req.FormValue("unitno")
	business.Website = req.FormValue("website")
	business.Instagram = req.FormValue("instagram")
	business.Facebook = req.FormValue("facebook")

	payload := struct {
		PageTitle  string
		User       *models.User
		Business   models.Business
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update Business", user, business, "", "",
	}

	data, err := json.Marshal(business)
	if err != nil {
		fmt.Println("error marshalling at process update business", err)
		payload.ErrorMsg = "An unexpected error has occured while updating business. Please try again."
		tpl.ExecuteTemplate(res, "updateBusiness.gohtml", payload)
		return
	}

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPut, apiBaseURL+"/businesses/"+vars["id"], bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lib.GetJWT(req))
	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending process update business request", err)
		payload.ErrorMsg = "An unexpected error has occured while updating business. Please try again."
	}
	defer response.Body.Close()

	// success
	if response.StatusCode == 200 {
		http.Redirect(res, req, "/business/"+vars["id"], http.StatusFound)
		return
	} else {
		// get error
		respData, _ := ioutil.ReadAll(response.Body)
		var errorResponse responses.ErrorResponse
		json.Unmarshal(respData, &errorResponse)
		payload.ErrorMsg = errorResponse.Error

		if strings.Contains(payload.ErrorMsg, "Duplicate entry") {
			payload.ErrorMsg = "There's already a business with the name provided. Please use a different name."
		}
	}

	tpl.ExecuteTemplate(res, "updateBusiness.gohtml", payload)
}

func ViewBusiness(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle   string
		User        *models.User
		Business    *models.Business
		Comments    []*models.Comment
		ErrorMsg    string
		SuccessMsg  string
		GMapsAPIKey string
	}{
		"View Business", nil, nil, nil, "", "", os.Getenv("GMapsPublicAPI"),
	}

	// Get User
	user, err := lib.IsLoggedIn(req)

	if err == nil {
		payload.User = user
	}

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, apiBaseURL+"/businesses/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending view business request", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		var business *models.Business
		marshalErr := json.Unmarshal(data, &business)

		if marshalErr != nil {
			fmt.Println("error unmasharling at view business", marshalErr)
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}

		payload.PageTitle = business.Name
		payload.Business = business

	} else {
		// handle error
		fmt.Println(string(data))
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	request, _ = http.NewRequest(http.MethodGet, apiBaseURL+"/business/"+vars["id"]+"/comments", nil)
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("error sending view business comments request", err)
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	data, _ = ioutil.ReadAll(response.Body)

	// success
	if response.StatusCode == 200 {
		var bComments []*models.Comment
		marshalErr := json.Unmarshal(data, &bComments)

		if marshalErr != nil {
			fmt.Println("error unmasharling at view business comments", marshalErr)
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}

		payload.Comments = bComments

		tpl.ExecuteTemplate(res, "viewBusiness.gohtml", payload)
		return

	} else {
		// handle error
		fmt.Println("Does this fail?")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}
}
