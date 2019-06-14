package main

type youTubePlugin struct {
}

func (y youTubePlugin) Name() string {
	return "Youtube Plugin"
}

// Plugin the plugin object to export
var Plugin youTubePlugin
