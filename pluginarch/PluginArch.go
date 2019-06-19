package pluginarch

type PluginConfig struct {
	Name    string
	Enabled bool
	Options map[string]string
}

// PluginAction Struct for plugin actions
type PluginAction struct {
	Name    string
	Type    string
	Options map[string]string
}
