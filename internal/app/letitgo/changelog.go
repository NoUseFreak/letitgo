package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/changelog"
)

func init() {
	registerAction(new(ChangelogAction))
}

type ChangelogAction struct{}

func (a *ChangelogAction) Name() string {
	return "changelog"
}

func (a *ChangelogAction) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"file": "CHANGELOG.md",
	}
}

func (a *ChangelogAction) GetDefaults() Config {
	return Config{
		Changelog: []changelog.Config{{
			File:    "CHANGELOG.md",
			Message: "Update changelog\n[skip ci]",
		}},
	}
}

func (a *ChangelogAction) Weight() int {
	return 5
}

func (a *ChangelogAction) Execute(cfg Config) error {
	for _, spec := range cfg.Changelog {
		spec.LetItGo = cfg.LetItGo
		if err := changelog.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
