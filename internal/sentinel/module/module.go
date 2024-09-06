package module

import (
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Module interface {
	Name() string
	Enable(s *discordgo.Session) error
	Disable(s *discordgo.Session) error
	Load(s *discordgo.Session) error
}

type ModuleManager struct {
	// NOTE: Only ever added to when initializing the bot,
	// there is no real way to dynamically register a module
	// that we did not have at build time. This is why we don't
	// need a lock here.
	Modules []Module
	mutex   *sync.RWMutex
}

func New() *ModuleManager {
	return &ModuleManager{}
}

func (m *ModuleManager) AddModule(module Module) {
	m.Modules = append(m.Modules, module)
}

// NOTE: This is where we enable all modules or disable them. This handles
// registration and removal of modules. Needs to be called once on initalization.
func (m *ModuleManager) LoadModules(s *discordgo.Session) {
	log.Info().Msg("Enabling modules...")
	for _, module := range m.Modules {
		err := module.Load(s)
		if err != nil {
			log.Error().Err(err).Str("module", module.Name()).Msg("Failed to load module")
		}
	}
	log.Info().Msg("Modules enabled")
}

func (m *ModuleManager) DisableByName(name string, s *discordgo.Session) {
	// Lock ourself up
	log.Info().Str("module", name).Msg("Enabling module")
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Find our culprit and disable it
	for _, module := range m.Modules {
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
	for _, module := range m.Modules {
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
