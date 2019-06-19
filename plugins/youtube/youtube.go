package main

import "github.com/martamius/AddaPlex/pluginarch"

type youTubePlugin struct {
}

func (y youTubePlugin) Name() string {
	return "Youtube"
}

func (y youTubePlugin) Identifier() string {
	return "youtube"
}

func (y youTubePlugin) ActionDefinitions() []pluginarch.PluginAction {

	s := make([]pluginarch.PluginAction, 2)
	s[0] = pluginarch.PluginAction{
		Name: "Download Video",
		Type: "url",
		Options: map[string]string{
			"regex": "youtube.com/watch",
		},
	}
	s[1] = pluginarch.PluginAction{
		Name: "Download Audio",
		Type: "url",
		Options: map[string]string{
			"regex": "youtube.com/watch",
		},
	}
	return s
}

func (y youTubePlugin) PerformAction(action string, options map[string]string) (string, bool) {
	switch action {
	case "Download Video":
		return "Ok, Ill download that video", true
	case "Download Audio":
		return "Ok, Ill download that audio", true

	default:
		return "Unkown action", false
	}
}

// Plugin the plugin object to export
var Plugin youTubePlugin
