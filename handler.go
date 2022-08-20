package main

import (
	"html/template"
	"log"
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
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseGlob("templates/*.html"))
	response, err := fetchData()
	if err != nil {
		log.Fatalf("Unable to load data from external API: %v", err)
	}
	data := IndexTemplateData{
		Author:   "Paul Burt",
		Email:    "paul.burt@bbc.co.uk",
		Schedule: response.Results,
		Title:    "Dagskrá RÚV",
		Today:    response.Results[0].StartTime.DateString(),
	}
	templateError := t.Execute(w, data)
	if templateError != nil {
		log.Fatalf("Unable to render HTML template: %v", templateError)
	}
}
