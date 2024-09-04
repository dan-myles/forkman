package module

import (
	"reflect"

	"github.com/avvo-na/devil-guard/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// All modules must implement this interface
type Module interface {
	Name() string
	Enable(s *discordgo.Session) error
	Disable(s *discordgo.Session) error
}

type ModuleManager struct {
	Modules []Module
}

func New() *ModuleManager {
	return &ModuleManager{}
}

func (m *ModuleManager) RegisterModule(module Module) {
	m.Modules = append(m.Modules, module)
}

func (m *ModuleManager) EnableModules(s *discordgo.Session) {
	cfg := config.GetConfig()
	cfg.RWMutex.RLock()
	defer cfg.RWMutex.RUnlock()

	for _, module := range m.Modules {
		// loop through cfg.ModuleCfg and enable the modules
		v := reflect.ValueOf(cfg.ModuleCfg)

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.IsValid() {
				// check if the module is enabled
				if field.Interface().(bool) {
					log.Info().Msgf("Module %s is enabled: %t", field.Type().Name(), field.Interface().(bool))
					module.Enable(s)
				}
			}
		}
	}

	// for _, module := range m.Modules {
	// 	module.Enable()
	// }
}

func (m *ModuleManager) DisableModules() {
	// for _, module := range m.Modules {
	// 	module.Disable()
	// }
}
