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
	"github.com/mofodox/project-live-app/client/lib"
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

	comment := req.FormValue("comment")
	var newComment models.Comment

	newComment.BusinessID = uint32(bID)
	newComment.Content = comment

	payload.Comment = &newComment

	data, err := json.Marshal(newComment)
	if err != nil {
		fmt.Println("error marshalling new comment", err)
	}
	reqBody := bytes.NewBuffer(data)

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/comment/", reqBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lib.GetJWT(req))

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
			http.Redirect(res, req, "/business/"+strconv.FormatUint(uint64(newComment.BusinessID), 10), http.StatusSeeOther)
			return
		}
	}
	tpl.ExecuteTemplate(res, "createComment.gohtml", payload)
}

func ViewComment(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["bID"])
	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}
	_, err = strconv.Atoi(vars["cID"])
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
		"View Comment", nil, "", "",
	}
	var comment models.Comment

	payload.Comment = &comment

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/comment/"+vars["cID"], nil)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Comment retrieval failed, error sending http request")
	} else {
		resBody, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == 201 {
			marshalErr := json.Unmarshal(resBody, &comment)
			if marshalErr != nil {
				payload.ErrorMsg = "An unexpected error has occured while showing business. Please try again."
				http.Redirect(res, req, "/business/"+vars["bID"], http.StatusSeeOther)
				return
			}
			fmt.Println("New comment created successfully")
			tpl.ExecuteTemplate(res, "viewComment.gohtml", payload)
			return
		}
	}
	tpl.ExecuteTemplate(res, "viewComment.gohtml", payload)
}

func UpdateComment(res http.ResponseWriter, req *http.Request) {
	payload := struct {
		PageTitle  string
		Comment    *models.Comment
		ErrorMsg   string
		SuccessMsg string
	}{
		"Edit Comment", nil, "", "",
	}
	tpl.ExecuteTemplate(res, "updateComment.gohtml", payload)

}

func ProcessUpdateComment(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["cID"])
	if err != nil {
		// Redirect to Index Page
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}
	_, err = strconv.Atoi(vars["bID"])
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
		"Edit Comment", nil, "", "",
	}
	var comment models.Comment
	content := req.FormValue("comment")
	comment.Content = content

	payload.Comment = &comment

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/comment/"+vars["cID"], nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lib.GetJWT(req))

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("")
		http.Redirect(res, req, "/business/"+vars["bID"], http.StatusSeeOther)
	} else {
		resBody, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == 200 {
			marshalErr := json.Unmarshal(resBody, &comment)

			if marshalErr != nil {
				fmt.Println("An unexpected error has occured while editing comment.")
				tpl.ExecuteTemplate(res, "updateComment.gohtml", payload)
			}
			fmt.Println("New comment created successfully")
			http.Redirect(res, req, "/business/"+strconv.FormatUint(uint64(comment.BusinessID), 10), http.StatusOK)
		}
	}
	tpl.ExecuteTemplate(res, "updateComment.gohtml", payload)
}

func DeleteComment(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	_, err := strconv.Atoi(vars["cID"])
	if err != nil {
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}
	_, err = strconv.Atoi(vars["bID"])
	if err != nil {
		http.Redirect(res, req, "/", http.StatusNotFound)
		return
	}
	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/comment/"+vars["cID"], nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+lib.GetJWT(req))

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Comment deletion failed, error sending http request")
		http.Redirect(res, req, "/business/"+vars["bID"], http.StatusSeeOther)
	} else {
		if response.StatusCode == 200 {
			fmt.Println("Comment deleted successfully")
			http.Redirect(res, req, "/business/"+vars["bID"], http.StatusSeeOther)
		}
	}
}
