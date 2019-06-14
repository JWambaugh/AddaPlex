package main

import (
	"log"

	"github.com/martamius/AddaPlex/pluginmanager"
)

func main() {
	log.Print("Starting AddaPlex")

	loadConfig()
	plugin := pluginmanager.LoadPlugin("youtube")
	log.Print("loaded plugin " + plugin.Name())

	startHTTP()
}
