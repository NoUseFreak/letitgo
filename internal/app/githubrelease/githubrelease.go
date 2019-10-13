package githubrelease

import (
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

func Execute(c Config) error {
	ui.Step("Publishing releases")

	templateProps(&c)

	if c.Owner == "" || c.Repo == "" {
		resolveOwnerRepo(&c)
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
		c.BaseConfig.LetItGo.Version(),
		c.Title,
		c.Description,
	)
	if err != nil {
		return err
	}

	return client.UploadAssets(rID, files)
}

func templateProps(c *Config) {
	utils.TemplateProperty(&c.Title, c)
	utils.TemplateProperty(&c.Description, c)
}

func resolveOwnerRepo(c *Config) error {
	repo, err := utils.GetRemote(".")
	if err != nil {
		return err
	}

	uri, err := utils.ParseURI(repo)
	if err != nil {
		return err
	}

	c.Owner = uri.Owner
	c.Repo = uri.Repo

	return nil
}
