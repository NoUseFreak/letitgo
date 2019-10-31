package homebrew

import (
	"fmt"
	"path"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/NoUseFreak/letitgo/internal/app/utils/git"
	"github.com/pkg/errors"
)

// New returns an action for homebrew
func New() action.Action {
	return &homebrew{}
}

type homebrew struct {
	Homepage string

	URL string

	Dependencies []string
	Conflicts    []string

	Tap     tapConfig
	Folder  string
	Install string
	Test    string
}

func (*homebrew) Name() string {
	return "homebrew"
}

func (*homebrew) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"homepage": "https://example.com",
		"url":      "https://github.com/owner/repo/releases/download/{{ .Version }}/darwin_amd64.zip",
		"tap": map[string]string{
			"url": "git@github.com:owner/homebrew-brew.git",
		},
	}
}

func (*homebrew) Weight() int {
	return 100
}

func (c *homebrew) Execute(cfg config.LetItGoConfig) error {
	setDefaults(c)
	templateProps(c, &cfg)
	hash, err := utils.BuildURLHash("sha256", c.URL)
	if err != nil {
		return errors.Wrapf(err, "Failed to build hash for %s", c.URL)
	}

	content, err := utils.Template(homebrewTpl, cfg, &c, map[string]string{
		"Hash": hash,
	})
	if err != nil {
		return err
	}

	ui.Trace(content)

	filename := path.Join(c.Folder, fmt.Sprintf("%s.rb", cfg.Name))
	message := fmt.Sprintf("Upgrade %s to %s", cfg.Name, cfg.Version())

	ui.Step("Uploading %s", filename)
	client, err := git.NewClient(
		c.Tap.URL,
		utils.DryRun.IsEnabled(),
	)
	if err != nil {
		return err
	}
	return client.PublishFile(filename, content, message, nil)
}
