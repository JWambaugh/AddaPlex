package main

import "github.com/martamius/AddaPlex/pluginarch"

type fooPlugin struct {
}

func (y fooPlugin) Name() string {
	return "Foo"
}

func (y fooPlugin) Identifier() string {
	return "foo"
}

func (y fooPlugin) ActionDefinitions() []pluginarch.PluginAction {
	s := make([]pluginarch.PluginAction, 1)
	s[0] = pluginarch.PluginAction{
		Name: "Say Foo",
		Type: "url",
	}
	return s
}
func (y fooPlugin) PerformAction(action string, options map[string]string) (string, bool) {
	return "FOO!" + options["url"], true
}

// Plugin the plugin object to export
var Plugin fooPlugin
