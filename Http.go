package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/martamius/AddaPlex/pluginarch"
	"github.com/martamius/AddaPlex/pluginmanager"
	"github.com/rs/cors"
)

type statusResponsePlugin struct {
	Name       string
	Identifier string
	Actions    []pluginarch.PluginAction
}
type statusResponse struct {
	Plugins    []statusResponsePlugin
	ServerName string
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	resp := statusResponse{}

	plugins := pluginmanager.Plugins()

	for _, plugin := range *plugins {
		p := statusResponsePlugin{}
		p.Name = plugin.Name()
		p.Identifier = plugin.Identifier()
		p.Actions = plugin.ActionDefinitions()

		resp.Plugins = append(resp.Plugins, p)
	}
	resp.ServerName = configData.ServerName
	js, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s", js)
}

type actionResponse struct {
	Status  string
	Message string
}

func performAction(w http.ResponseWriter, r *http.Request) {
	plugins := pluginmanager.Plugins()
	for _, plugin := range *plugins {
		if plugin.Identifier() == r.FormValue("pluginIdentifier") {
			response := actionResponse{}
			var ok bool
			options := make(map[string]string)
			options["url"] = r.FormValue("url")
			response.Message, ok = plugin.PerformAction(r.FormValue("Name"), options)
			if ok {
				response.Status = "ok"
			} else {
				response.Status = "error"
			}
			js, _ := json.Marshal(response)
			fmt.Fprintf(w, "%s", js)
			return
		}
	}
	response := actionResponse{}
	response.Message = "Invalid plugin identifier"
	response.Status = "error"
	js, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", js)
}

func startHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	mux.HandleFunc("/status", getStatus)
	mux.HandleFunc("/action", performAction)

	log.Print("Listening on port " + configData.ListenPort + " as " + configData.ServerName)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":"+configData.ListenPort, handler))
}
