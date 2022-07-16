package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type CustomTime time.Time

const timeFormat = "2006-01-02 15:04:05"

func (c *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		*c = CustomTime(t)
	}
	return
}

func (c *CustomTime) DateString() string {
	t := time.Time(*c)
	return t.Format("02.01.2006")
}

func (c *CustomTime) TimeString() string {
	t := time.Time(*c)
	return t.Format("15:04")
}

type Response struct {
	Results []Show
}

type Show struct {
	IsLive              bool       `json:"live"`
	OriginalDescription string     `json:"description"`
	StartTime           CustomTime `json:"startTime"`
	Title               string     `json:"title"`
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
	return s.StartTime.TimeString()
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
	t := template.Must(template.ParseFiles("templates/index.gohtml"))
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
	t := template.Must(template.ParseFiles("templates/schedule.gohtml"))
	data := ScheduleData{Schedule: getSchedule()}
	err := t.Execute(w, data)
	if err != nil {
		log.Fatalf("Unable to render HTML template: %s", err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/schedule", ScheduleHandler).Methods("GET").Headers("HX-Request", "")
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", r)
}
