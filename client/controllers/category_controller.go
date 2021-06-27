package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
)

func ViewCategory(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	// anonymous payload
	payload := struct {
		PageTitle  string
		Category   *models.Category
		ErrorMsg   string
		SuccessMsg string
	}{
		"View Category", nil, "", "",
	}

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/categories/"+vars["id"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	// handle error
	if err != nil {
		fmt.Println("Fetch category at view category failed")
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// success
		if response.StatusCode == 200 {
			var category *models.Category
			marshalErr := json.Unmarshal(data, &category)

			if marshalErr != nil {
				fmt.Println("Error decoding json at view category")
				http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
				return
			}

			payload.PageTitle = category.Name
			payload.Category = category

			tpl.ExecuteTemplate(res, "viewCategory.gohtml", payload)
			return
		} else {
			// handle error
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
			return
		}
	}
}

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
	fmt.Println(data)

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/categories", respBody)
	fmt.Println(request)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	fmt.Println(response)

	if err != nil {
		// handle error
		fmt.Println("Category creation failed")
		http.Redirect(res, req, "/", response.StatusCode)
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
			http.Redirect(res, req, "/", response.StatusCode)
		}
	}
}

func UpdateCategory(res http.ResponseWriter, req *http.Request) {

	fmt.Println("UPDATE CATEGORY FORM")

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	var category *models.Category

	// Anonymous payload
	payload := struct {
		PageTitle  string
		Category   *models.Category
		ErrorMsg   string
		SuccessMsg string
	}{
		"Update Category", nil, "", "",
	}

	// Todo: add cookie check and send JWT with request
	client := &http.Client{}
	fmt.Println(client)
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/categories/"+vars["id"], nil)
	fmt.Println(request)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	fmt.Println(response)

	if err != nil {
		// handle error
		fmt.Println("Category update failed")
		http.Redirect(res, req, "/", response.StatusCode)
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		fmt.Println(response.StatusCode)

		data, _ := ioutil.ReadAll(response.Body)
		marshalErr := json.Unmarshal(data, &category)

		if marshalErr != nil {
			fmt.Println("Error decoding json")
		}

		payload.Category = category

		tpl.ExecuteTemplate(res, "updateCategory.gohtml", payload)
		return
	}
}

func ProcessUpdateCategory(res http.ResponseWriter, req *http.Request) {

	fmt.Println("PROCESS UPDATE CATEGORY")

	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

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
	request, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/categories/"+vars["id"], respBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		// handle error
		fmt.Println("Category creation failed")
		http.Redirect(res, req, "/", response.StatusCode)
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)

		// Success
		if response.StatusCode == 200 {
			fmt.Println("Category updated successfully")
			fmt.Println(string(data))
			http.Redirect(res, req, "/", http.StatusOK)
		} else {
			// handle error
			fmt.Println("Category update failed")
			http.Redirect(res, req, "/", response.StatusCode)
		}
	}
}
