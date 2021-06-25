package controllers

import (
	"log"
	"net/http"
	"text/template"
)

func Home(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/index.gohtml")
	
	if err != nil {
		log.Fatalf("Error template home %v\n", err)
	}

	t.Execute(res, nil)
}
