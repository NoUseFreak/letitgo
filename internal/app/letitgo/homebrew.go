package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
)

func init() {
	RegisterAction(new(HomebrewAction))
}

type HomebrewAction struct{}

func (a *HomebrewAction) Weight() int {
	return 100
}

func (a *HomebrewAction) Execute(cfg Config) error {
	for _, spec := range cfg.Homebrew {
		spec.Name = cfg.Name
		spec.Version = cfg.Version
		spec.Description = cfg.Description
		if err := homebrew.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
