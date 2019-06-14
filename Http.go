package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/martamius/AddaPlex/pluginarch"
	"github.com/martamius/AddaPlex/pluginmanager"
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

func startHTTP() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/status", getStatus)

	log.Print("Listening on port " + configData.ListenPort + " as " + configData.ServerName)

	log.Fatal(http.ListenAndServe(":"+configData.ListenPort, nil))
}
