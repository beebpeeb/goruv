package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func fetchData() (r Response, err error) {
	var c = &http.Client{
		Timeout: time.Second * 10,
	}
	res, err := c.Get("https://apis.is/tv/ruv")
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	var data Response
	jsonError := json.Unmarshal(body, &data)
	if jsonError != nil {
		return
	}
	return data, err
}
