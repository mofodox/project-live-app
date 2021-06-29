package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/api/responses"
)

func CreateComment(res http.ResponseWriter, req *http.Request) {

}

func ProcessCommentForm(res http.ResponseWriter, req *http.Request) {
	businessInfo := req.FormValue("business ID")
	comment := req.FormValue("content")
	var newComment models.Comment
	bID, err := strconv.Atoi(businessInfo)
	if err != nil {
		responses.ERROR(res, http.StatusInternalServerError, err)
		return
	}
	newComment.BusinessID = uint32(bID)
	newComment.Content = comment

	data, err := json.Marshal(newComment)
	if err != nil {
		fmt.Println("error marshalling new comment", err)
	}
	resBody := bytes.NewBuffer(data)

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/comment/", resBody)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Comment creation failed, error sending http request")
		http.Redirect(res, req, "/", response.StatusCode)
	}
}
