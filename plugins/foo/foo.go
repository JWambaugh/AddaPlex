package main

import "github.com/martamius/AddaPlex/pluginarch"

type fooPlugin struct {
}

func (y fooPlugin) Name() string {
	return "Foo - I just say foo!"
}

func (y fooPlugin) Identifier() string {
	return "foo"
}

func (y fooPlugin) ActionDefinitions() []pluginarch.PluginAction {
	a := pluginarch.PluginAction{}
	s := make([]pluginarch.PluginAction, 1)
	s[0] = a
	return s
}

// Plugin the plugin object to export
var Plugin fooPlugin
