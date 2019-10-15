package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/snapcraft"
)

func init() {
	registerAction(new(SnapcraftAction))
}

type SnapcraftAction struct{}

func (a *SnapcraftAction) Name() string {
	return "snapcraft"
}

func (a *SnapcraftAction) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"assets": []string{
			"./build/bin/linux_amd64/letitgo",
		},
		"architecture": "amd64",
	}
}

func (a *SnapcraftAction) GetDefaults() Config {
	return Config{
		Snapcraft: []snapcraft.Config{{}},
	}
}

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
