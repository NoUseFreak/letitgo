package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
)

func init() {
	registerAction(new(HomebrewAction))
}

type HomebrewAction struct{}

func (a *HomebrewAction) Name() string {
	return "homebrew"
}

func (a *HomebrewAction) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"homepage": "https://example.com",
		"url":      "https://github.com/owner/repo/releases/download/{{ .Version }}/darwin_amd64.zip",
		"tap": map[string]string{
			"url": "git@github.com:owner/homebrew-brew.git",
		},
	}
}

func (a *HomebrewAction) GetDefaults() Config {
	return Config{
		Homebrew: []homebrew.Config{{}},
	}
}

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
