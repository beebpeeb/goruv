package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Results []Show
}

type Show struct {
	IsLive              bool   `json:"live"`
	OriginalDescription string `json:"description"`
	StartTime           string `json:"startTime"`
	Title               string `json:"title"`
}

func (s Show) Description() string {
	return strings.TrimSuffix(s.OriginalDescription, " e.")
}

func (s Show) HasDescription() bool {
	return len(strings.TrimSpace(s.OriginalDescription)) >= 1
}

func (s Show) IsRepeat() bool {
	return strings.HasSuffix(s.OriginalDescription, " e.")
}

func (s Show) Time() string {
	t, err := time.Parse("2006-01-02 15:04:05", s.StartTime)
	if err != nil {
		log.Fatalf("Unable to parse date/time string: %s", err)
	}
	return t.Format("15:04")
}

func getSchedule() []Show {
	url := "https://apis.is/tv/ruv"
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %s", url, err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Unable to read response body: %s", err)
	}
	var r Response
	jsonErr := json.Unmarshal(body, &r)
	if jsonErr != nil {
		log.Fatalf("Unable to parse JSON: %s", jsonErr)
	}
	return r.Results
}

type IndexData struct {
	Author string
	Email  string
	Title  string
	Today  string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	data := IndexData{
		Author: "Paul Burt",
		Email:  "paul.burt@bbc.co.uk",
		Title:  "Dagskrá RÚV",
		Today:  time.Now().Format("02.01.2006"),
	}
	err := t.Execute(w, data)
	if err != nil {
		log.Fatalf("Unable to render HTML template: %s", err)
	}
}

type ScheduleData struct {
	Schedule []Show
}

func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Header.Get("HX-Request")) == "true" {
		t := template.Must(template.ParseFiles("templates/schedule.html"))
		data := ScheduleData{Schedule: getSchedule()}
		err := t.Execute(w, data)
		if err != nil {
			log.Fatalf("Unable to render HTML template: %s", err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	// Handlers
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/schedule", ScheduleHandler)

	// Static files
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Server
	http.ListenAndServe(":8080", nil)
}
