package action

import "github.com/NoUseFreak/letitgo/internal/app/config"

// Action is the action interface
type Action interface {
	Name() string
	GetInitConfig() map[string]interface{}
	Weight() int
	Execute(config.LetItGoConfig) error
}
