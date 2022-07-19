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

type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		*ct = CustomTime(t)
	}
	return err
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

func getSchedule() ([]Show, error) {
	res, httpErr := http.Get("https://apis.is/tv/ruv")
	if httpErr != nil {
		return nil, httpErr
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	var r Response
	unmarshalErr := json.Unmarshal(body, &r)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return r.Results, nil
}

type IndexTemplateData struct {
	Author   string
	Email    string
	Schedule []Show
	Title    string
	Today    string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseGlob("templates/*.html"))
	schedule, err := getSchedule()
	if err != nil || len(schedule) == 0 {
		log.Fatalf("Unable to load data from external API: %s", err)
	}
	data := IndexTemplateData{
		Author:   "Paul Burt",
		Email:    "paul.burt@bbc.co.uk",
		Schedule: schedule,
		Title:    "Dagskrá RÚV",
		Today:    schedule[0].StartTime.DateString(),
	}
	templateErr := t.Execute(w, data)
	if templateErr != nil {
		log.Fatalf("Unable to render HTML template: %s", templateErr)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", mux)
}
