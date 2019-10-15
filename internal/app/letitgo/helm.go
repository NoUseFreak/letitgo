package letitgo

import "github.com/NoUseFreak/letitgo/internal/app/helm"

func init() {
	registerAction(new(HelmAction))
}

type HelmAction struct{}

func (a *HelmAction) Name() string {
	return "helm"
}

func (a *HelmAction) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"charts":     []string{"./charts/my-chart"},
		"builddir":   "./build/helm",
		"repository": "https://chartmuseum.example.com",
	}
}

func (a *HelmAction) GetDefaults() Config {
	return Config{
		Helm: []helm.Config{{}},
	}
}

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
