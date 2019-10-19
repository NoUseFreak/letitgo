package githubrelease

import (
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

// Action creates a release and attaches any given artifacts.
type Action struct {
	Owner string
	Repo  string

	Title       string
	Description string

	Assets []assetConfig
}

// Name return the name of the action.
func (*Action) Name() string {
	return "githubrelease"
}

// GetInitConfig return what a good starting config would be.
func (*Action) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"assets": []string{"./build/*"},
	}
}

// Weight return in what order this action should be handled.
func (*Action) Weight() int {
	return 10
}

// Execute handles the action.
func (c *Action) Execute(cfg config.LetItGoConfig) error {
	templateProps(c, &cfg)

	if c.Owner == "" || c.Repo == "" {
		resolveOwnerRepo(c)
	}

	files := []string{}
	for _, a := range c.Assets {
		for _, f := range a.GetFiles() {
			files = append(files, f)
		}
	}

	client := utils.GithubClient{
		Owner: c.Owner,
		Repo:  c.Repo,
	}
	rID, err := client.CreateRelease(
		cfg.Version(),
		c.Title,
		c.Description,
	)
	if err != nil {
		return err
	}

	return client.UploadAssets(rID, files)
}
