package githubrelease

import (
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	gitclient "github.com/NoUseFreak/letitgo/internal/app/utils/git"
)

func templateProps(c *githubrelease, cfg *config.LetItGoConfig) {
	utils.TemplateProperty(&c.Title, c, cfg)
	utils.TemplateProperty(&c.Description, c, cfg)
}

func resolveOwnerRepo(c *githubrelease) error {
	repo, err := gitclient.GetRemote(".")
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
