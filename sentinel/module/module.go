package module

import (
	"reflect"
	"strings"

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

func (m *ModuleManager) AddModule(module Module) {
	m.Modules = append(m.Modules, module)
}

func (m *ModuleManager) RegisterModules(s *discordgo.Session) {
	log.Info().Msg("Enabling modules...")
	cfg := config.GetConfig()
	cfg.RWMutex.RLock()
	defer cfg.RWMutex.RUnlock()

	// Loop through all modules
	// Check if the module is enabled in config
	// Then enable it 😊
	for _, module := range m.Modules {
		// loop through cfg.ModuleCfg and enable the modules
		v := reflect.ValueOf(cfg.ModuleCfg)
		log.Debug().Str("module", module.Name()).Msg("Found interfaced module in register")

		for i := 0; i < v.NumField(); i++ {
			// Get field and check if it's valid
			field := v.Field(i)
			if !field.IsValid() {
				log.Debug().Str("module", module.Name()).Msg("Module field is not valid")
				continue
			}

			// Check if the module is enabled
			enabled := field.FieldByName("Enabled")
			if !enabled.IsValid() || enabled.Kind() != reflect.Ptr {
				log.Debug().Str("module", module.Name()).Msg("Module enabled field is not a pointer")
				continue
			}

			// Check if the value is a boolean
			enabledValue := enabled.Elem()
			if !enabledValue.IsValid() || enabledValue.Kind() != reflect.Bool {
				log.Debug().Str("module", module.Name()).Msg("Module enabled value is not a boolean")
				continue
			}

			// Check if the module name matches the field name
			if !strings.EqualFold(module.Name(), v.Type().Field(i).Name) {
				log.Debug().Str("module", module.Name()).Str("field", v.Type().Field(i).Name).Msg("Module name does not match field name")
				continue
			}

			// INFO: Here is where we actually enable/disable the module
			// Now we either enable or disable the module!
			if enabledValue.Bool() {
				err := module.Enable(s)
				if err != nil {
					log.Error().Err(err).Str("module", module.Name()).Msg("Failed to enable module")
				}

				log.Info().Str("module", module.Name()).Msg("Module enabled")
				break
			} else {
				err := module.Disable(s)
				if err != nil {
					log.Error().Err(err).Str("module", module.Name()).Msg("Failed to disable module")
				}

				log.Info().Str("module", module.Name()).Msg("Module disabled")
				break
			}
		}
	}
}
