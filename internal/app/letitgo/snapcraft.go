package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/snapcraft"
)

func init() {
	registerAction(new(SnapcraftAction))
}

type SnapcraftAction struct{}

func (a *SnapcraftAction) Weight() int {
	return 110
}

func (a *SnapcraftAction) Execute(cfg Config) error {
	for _, spec := range cfg.Snapcraft {
		spec.LetItGo = cfg.LetItGo
		if err := snapcraft.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
