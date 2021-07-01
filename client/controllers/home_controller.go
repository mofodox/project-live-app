package controllers

import (
	"net/http"

	"github.com/mofodox/project-live-app/api/models"
	"github.com/mofodox/project-live-app/client/lib"
)

func Home(res http.ResponseWriter, req *http.Request) {

	// anonymous payload
	payload := struct {
		PageTitle  string
		User       *models.User
		ErrorMsg   string
		SuccessMsg string
	}{
		"Businesses", nil, "", "",
	}

	// Get User
	user, err := lib.IsLoggedIn(req)

	if err == nil {
		payload.User = user
	}

	lib.Tpl.ExecuteTemplate(res, "index.gohtml", payload)
}
