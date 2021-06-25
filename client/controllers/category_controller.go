package controllers

import (
	"fmt"
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
