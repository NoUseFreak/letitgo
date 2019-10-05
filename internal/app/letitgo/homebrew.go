package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
)

func init() {
	registerAction(new(HomebrewAction))
}

type HomebrewAction struct{}

func (a *HomebrewAction) Weight() int {
	return 100
}

func (a *HomebrewAction) Execute(cfg Config) error {
	for _, spec := range cfg.Homebrew {
		spec.LetItGo = cfg.LetItGo
		if err := homebrew.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
