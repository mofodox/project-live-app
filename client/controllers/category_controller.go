package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mofodox/project-live-app/api/models"
)

func CreateCategoryPage(res http.ResponseWriter, req *http.Request) {

	fmt.Println("CREATE Category Page")

	// Anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.Category
		ErrorMsg   string
		SuccessMsg string
	}{
		"Create Category", nil, "", "",
	}

	tpl.ExecuteTemplate(res, "createCategory.gohtml", payload)
}

func ProcessCategoryForm(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Post Category Page")
	categoryName := req.FormValue("name")
	description := req.FormValue("description")

	data, err := json.Marshal(map[string]string{
		"name":        categoryName,
		"description": description,
	})
	if err != nil {
		log.Fatalf("login error %v\n", err)
	}

	respBody := bytes.NewBuffer(data)

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/categories", respBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)

		// Success
		if response.StatusCode == 201 {
			fmt.Println("Category created successfully")
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusOK)
		} else {
			// handle error
			fmt.Println("Category creation failed")
		}
	}
}

// func UpdateCategoryPage(res http.ResponseWriter, req *http.Request) {

// 	fmt.Println("UPDATE BUSINESS FORM")

// 	vars := mux.Vars(req)
// 	category_id, err := strconv.Atoi(vars["id"])

// 	if err != nil {
// 		// Redirect to Index Page
// 		http.Redirect(res, req, "/", http.StatusNotFound)
// 		return
// 	}

// 	var category *models.Category

// 	if err := server.DB.First(&category, category_id).Error; err != nil {
// 		// Redirect to Index Page
// 		http.Redirect(res, req, "/", http.StatusNotFound)
// 		return
// 	}

// 	fmt.Println(category)

// 	// Anonymous payload
// 	payload := struct {
// 		PageTitle  string
// 		User       *models.User
// 		Business   *models.Category
// 		ErrorMsg   string
// 		SuccessMsg string
// 	}{
// 		"Update Category", nil, category, "", "",
// 	}

// 	tpl.ExecuteTemplate(res, "updateCategory.gohtml", payload)
// }
