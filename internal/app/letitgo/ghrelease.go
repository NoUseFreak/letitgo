package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/githubrelease"
)

func init() {
	registerAction(new(GithubReleaseAction))
}

type GithubReleaseAction struct{}

func (a *GithubReleaseAction) Name() string {
	return "githubrelease"
}

func (a *GithubReleaseAction) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"assets": []string{"./build/*"},
	}
}

func (a *GithubReleaseAction) GetDefaults() Config {
	return Config{
		GithubRelease: []githubrelease.Config{{}},
	}
}

func (a *GithubReleaseAction) Weight() int {
	return 10
}

func (a *GithubReleaseAction) Execute(cfg Config) error {
	for _, spec := range cfg.GithubRelease {
		spec.LetItGo = cfg.LetItGo
		if err := githubrelease.Execute(spec); err != nil {
			return err
		}
	}

	return nil
}
