package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mofodox/project-live-app/api/models"
)

func CreateComment(res http.ResponseWriter, req *http.Request) {
	// Get User
	/*user, err := IsLoggedIn(req)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusSeeOther)
		return
	} */

	payload := struct {
		PageTitle string
		//User       *models.User
		Comment    *models.Comment
		ErrorMsg   string
		SuccessMsg string
	}{
		"New Comment" /*user,*/, nil, "", "",
	}

	tpl.ExecuteTemplate(res, "createComment.gohtml", payload)
}

func ProcessCommentForm(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	bID, err := strconv.Atoi(vars["id"])

	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}

	payload := struct {
		PageTitle  string
		Comment    *models.Comment
		ErrorMsg   string
		SuccessMsg string
	}{
		"Create Business", nil, "", "",
	}

	comment := req.FormValue("content")
	var newComment models.Comment

	newComment.BusinessID = uint32(bID)
	newComment.Content = comment

	data, err := json.Marshal(newComment)
	if err != nil {
		fmt.Println("error marshalling new comment", err)
	}
	reqBody := bytes.NewBuffer(data)

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/comment/", reqBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Comment creation failed, error sending http request")
		http.Redirect(res, req, "/", response.StatusCode)
	} else {
		resBody, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == 201 {
			marshalErr := json.Unmarshal(resBody, &newComment)
			if marshalErr != nil {
				payload.ErrorMsg = "An unexpected error has occured while creating business. Please try again."
				tpl.ExecuteTemplate(res, "createComment.gohtml", payload)
				return
			}
			fmt.Println("New comment created successfully")
			http.Redirect(res, req, "/business/"+strconv.FormatUint(uint64(newComment.BusinessID), 10), http.StatusFound)
			return
		}
	}
	tpl.ExecuteTemplate(res, "createComment.gohtml", payload)
}
