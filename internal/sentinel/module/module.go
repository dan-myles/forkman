package module

import (
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Module interface {
	// Returns the name of the module
	Name() string

	// Enables the module, handles any setup and registration
	// of commands, writes config to file.
	Enable(s *discordgo.Session) error

	// Disables the module, handles any cleanup and deregistration
	// of commands, writes config to file.
	Disable(s *discordgo.Session) error

	// Loads the module, handles any setup and registration of
	// commands, *reads* config from file. To only be called once
	Load(s *discordgo.Session) error
}

type ModuleManager struct {
	modules []Module
	mutex   sync.RWMutex
}

func New() *ModuleManager {
	return &ModuleManager{}
}

// NOTE:
// All modules should be added to the manager
// regardless of whether they are enabled or not.
func (m *ModuleManager) AddModule(module Module) {
	// Lock ourself up
	log.Info().
		Str("module", module.Name()).
		Msg("Adding module to be loaded...")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.modules = append(m.modules, module)
}

func (m *ModuleManager) LoadModules(s *discordgo.Session) {
	// Lock ourself up
	log.Debug().Msg("Loading modules...")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, module := range m.modules {
		err := module.Load(s)
		if err != nil {
			log.Error().Err(err).Str("module", module.Name()).Msg("Failed to load module")
		}
	}
}

func (m *ModuleManager) DisableByName(name string, s *discordgo.Session) {
	// Lock ourself up
	log.Info().Str("module", name).Msg("Enabling module")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Find our culprit and disable it
	for _, module := range m.modules {
		if strings.EqualFold(module.Name(), name) {
			err := module.Disable(s)
			if err != nil {
				log.Error().Err(err).Str("module", module.Name()).Msg("Failed to disable module")
			}

			log.Info().Str("module", module.Name()).Msg("Module disabled")
			break
		}
	}
}

func (m *ModuleManager) EnableByName(name string, s *discordgo.Session) {
	// Lock ourself up
	log.Info().Str("module", name).Msg("Disabling module")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Find our culprit and enable it
	for _, module := range m.modules {
		if strings.EqualFold(module.Name(), name) {
			err := module.Enable(s)
			if err != nil {
				log.Error().Err(err).Str("module", module.Name()).Msg("Failed to enable module")
			}

			log.Info().Str("module", module.Name()).Msg("Module enabled")
			break
		}
	}
}
