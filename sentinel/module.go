package sentinel

type Module interface {
	Name() string
	Enable() error
	Disable() error
}

type ModuleManager struct {
	Modules []Module
}

func (m *ModuleManager) RegisterModule(module Module) {
	m.Modules = append(m.Modules, module)
}

func (m *ModuleManager) EnableModules() {
	for _, module := range m.Modules {
		module.Enable()
	}
}

func (m *ModuleManager) DisableModules() {
	for _, module := range m.Modules {
		module.Disable()
	}
}
