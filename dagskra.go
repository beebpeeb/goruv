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

func (ct *CustomTime) DateString() string {
	t := time.Time(*ct)
	return t.Format("02.01.2006")
}

func (c *CustomTime) TimeString() string {
	t := time.Time(*c)
	return t.Format("15:04")
}

type Response struct {
	Results []Show `json:"results"`
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

func getResponse() (Response, error) {
	res, err := http.Get("https://apis.is/tv/ruv")
	emptyResponse := Response{Results: []Show{}}
	if err != nil {
		return emptyResponse, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return emptyResponse, err
	}
	var r Response
	jsonErr := json.Unmarshal(body, &r)
	if jsonErr != nil {
		return emptyResponse, jsonErr
	}
	return r, nil
}

type IndexTemplateData struct {
	Author   string
	Email    string
	Error    error
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
	response, err := getResponse()
	data := IndexTemplateData{
		Author:   "Paul Burt",
		Email:    "paul.burt@bbc.co.uk",
		Error:    err,
		Schedule: response.Results,
		Title:    "Dagskrá RÚV",
		Today:    response.Results[0].StartTime.DateString(),
	}
	tplErr := t.Execute(w, data)
	if tplErr != nil {
		log.Fatalf("Unable to render HTML template: %s", tplErr)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", mux)
}
