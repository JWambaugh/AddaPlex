package main

import (
	"log"

	"github.com/martamius/AddaPlex/pluginmanager"
)

func main() {
	log.Print("Starting AddaPlex")

	loadConfig()

	pluginmanager.LoadPlugins(moduleNames())

	startHTTP()
}
