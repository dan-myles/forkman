package discord

type Module interface {
	Load() error
	Enable() error
	Disable() error
	DisableCommand(string) error
	EnableCommand(string) error
}
