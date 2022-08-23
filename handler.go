package main

import (
	"html/template"
	"net/http"
)

type IndexTemplateData struct {
	Author   string
	Email    string
	Schedule []Listing
	Title    string
	Today    string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("templates/*.html"))
	response, _ := fetchData()
	data := IndexTemplateData{
		Author:   "Paul Burt",
		Email:    "paul.burt@bbc.co.uk",
		Schedule: response.Results,
		Title:    "Dagskrá RÚV",
		Today:    response.Results[0].StartTime.DateString(),
	}
	t.Execute(w, data)
}
