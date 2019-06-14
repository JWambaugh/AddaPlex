package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ModuleConfig struct {
	Name    string
	Enabled bool
}
type Config struct {
	ServerName string `yaml:"serverName"`
	ListenPort string `yaml:"listenPort"`
	Modules    []ModuleConfig
}

var configData Config

func moduleNames() []string {
	var names []string
	for _, module := range configData.Modules {
		if module.Name == "Other" {
			continue
		}
		if module.Enabled {
			names = append(names, module.Name)
		}

	}
	return names
}

func loadConfig() {
	configData = Config{}

	file, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	yml, err := ioutil.ReadAll(file)

	err = yaml.Unmarshal(yml, &configData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("--- t:\n%v\n\n", configData)

	module := ModuleConfig{Name: "Other"}
	configData.Modules = append(configData.Modules, module)

	// d, err := yaml.Marshal(&configData)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t dump:\n%s\n\n", string(d))
	log.Print("Loaded config")
}
