package pluginmanager

import (
	"fmt"
	"os"
	"plugin"
)

// AddaplexPlugin interface for plugins
type AddaplexPlugin interface {
	Name() string
}

var plugins []*AddaplexPlugin

// LoadPlugin Loads a plugin
func LoadPlugin(pluginName string) AddaplexPlugin {
	// determine plugin to load

	var mod string

	mod = "plugins/" + pluginName + ".so"

	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Plugin
	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type Greeter (defined above)
	var plugin AddaplexPlugin
	plugin, ok := symPlugin.(AddaplexPlugin)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}
	plugins = append(plugins, &plugin)

	return plugin
}
