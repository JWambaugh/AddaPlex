package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/martamius/AddaPlex/pluginarch"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerName string `yaml:"serverName"`
	ListenPort string `yaml:"listenPort"`
	Modules    []pluginarch.PluginConfig
}

var configData Config

func pluginConfigs() []pluginarch.PluginConfig {
	var plugins []pluginarch.PluginConfig
	for _, module := range configData.Modules {
		if module.Name == "Other" {
			continue
		}
		if module.Enabled {
			plugins = append(plugins, module)
		}

	}
	return plugins
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

	module := pluginarch.PluginConfig{Name: "Other"}
	configData.Modules = append(configData.Modules, module)

	// d, err := yaml.Marshal(&configData)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t dump:\n%s\n\n", string(d))
	log.Print("Loaded config")
}
