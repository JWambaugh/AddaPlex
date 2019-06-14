package main

import "github.com/martamius/AddaPlex/pluginarch"

type youTubePlugin struct {
}

func (y youTubePlugin) Name() string {
	return "Youtube Plugin"
}

func (y youTubePlugin) Identifier() string {
	return "youtube"
}

func (y youTubePlugin) ActionDefinitions() []pluginarch.PluginAction {

	s := make([]pluginarch.PluginAction, 1)
	s[0] = pluginarch.PluginAction{
		Name: "Download Video",
		Type: "url",
		Options: map[string]string{
			"regex": "youtoube.com",
		},
	}
	return s
}

// Plugin the plugin object to export
var Plugin youTubePlugin
