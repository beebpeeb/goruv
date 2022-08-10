package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", mux)
}
