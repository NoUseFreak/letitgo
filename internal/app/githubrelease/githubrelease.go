package githubrelease

import (
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

func templateProps(c *Action, cfg *config.LetItGoConfig) {
	utils.TemplateProperty(&c.Title, c, cfg)
	utils.TemplateProperty(&c.Description, c, cfg)
}

func resolveOwnerRepo(c *Action) error {
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
