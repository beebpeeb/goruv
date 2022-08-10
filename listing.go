package main

import (
	"strings"
	"time"
)

type customTime time.Time

func (ct *customTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		*ct = customTime(t)
	}
	return err
}

func (ct *customTime) DateString() string {
	t := time.Time(*ct)
	return t.Format("02.01.2006")
}

func (ct *customTime) TimeString() string {
	t := time.Time(*ct)
	return t.Format("15:04")
}

type Response struct {
	Results []Listing `json:"results"`
}

type Listing struct {
	IsLive              bool       `json:"live"`
	OriginalDescription string     `json:"description"`
	StartTime           customTime `json:"startTime"`
	Title               string     `json:"title"`
}

func (s *Listing) Description() string {
	return strings.TrimSuffix(s.OriginalDescription, " e.")
}

func (s *Listing) HasDescription() bool {
	return strings.TrimSpace(s.OriginalDescription) != ""
}

func (s *Listing) IsRepeat() bool {
	return strings.HasSuffix(s.OriginalDescription, " e.")
}

func (s *Listing) Time() string {
	return s.StartTime.TimeString()
}
