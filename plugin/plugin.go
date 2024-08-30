package plugin

type Plugin interface {
	// Name returns the name of the plugin
	Name() string
	// Register is used to register the plugin
	Register()
	// Run is the main function of the plugin
	Run() error
}

type PluginManager struct {
	Plugins []Plugin
}

func New() *PluginManager {
	return &PluginManager{}
}

func (pm *PluginManager) Register(p Plugin) {
	pm.Plugins = append(pm.Plugins, p)
}

func (pm *PluginManager) RegisterAll() {
	for _, p := range pm.Plugins {
		p.Register()
	}
}

func (pm *PluginManager) Run() {
	for _, p := range pm.Plugins {
		p.Run()
	}
}
