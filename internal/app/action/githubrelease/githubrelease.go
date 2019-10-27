package githubrelease

import (
	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

// New returns an action for githubrelease
func New() action.Action {
	return &githubrelease{}
}

type githubrelease struct {
	Owner string
	Repo  string

	Title       string
	Description string

	Assets []assetConfig
}

func (*githubrelease) Name() string {
	return "githubrelease"
}

func (*githubrelease) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"assets": []string{"./build/*"},
	}
}

func (*githubrelease) Weight() int {
	return 20
}

func (c *githubrelease) Execute(cfg config.LetItGoConfig) error {
	templateProps(c, &cfg)

	if c.Owner == "" || c.Repo == "" {
		if err := resolveOwnerRepo(c); err != nil {
			return err
		}
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
		cfg.Version().String(),
		c.Title,
		c.Description,
	)
	if err != nil {
		return err
	}

	return client.UploadAssets(rID, files)
}
