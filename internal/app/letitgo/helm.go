package letitgo

import "github.com/NoUseFreak/letitgo/internal/app/helm"

func init() {
	registerAction(new(HelmAction))
}

type HelmAction struct{}

func (a *HelmAction) Weight() int {
	return 120
}

func (a *HelmAction) Execute(cfg Config) error {
	for _, spec := range cfg.Helm {
		spec.LetItGo = cfg.LetItGo
		if err := helm.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
