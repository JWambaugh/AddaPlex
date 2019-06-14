package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func startHTTP() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Print("Listening on port " + configData.ListenPort + " as " + configData.ServerName)

	log.Fatal(http.ListenAndServe(":"+configData.ListenPort, nil))
}
