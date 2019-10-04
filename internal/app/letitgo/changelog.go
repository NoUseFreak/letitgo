package letitgo

import "github.com/NoUseFreak/letitgo/internal/app/changelog"

func init() {
	RegisterAction(new(ChangelogAction))
}

type ChangelogAction struct{}

func (a *ChangelogAction) Weight() int {
	return 5
}

func (a *ChangelogAction) Execute(cfg Config) error {
	for _, spec := range cfg.Changelog {
		if err := changelog.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
