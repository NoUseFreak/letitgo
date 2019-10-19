package homebrew

import (
	"fmt"
	"path"

	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

// Action creates or updates a homebrew tap.
type Action struct {
	Homepage string

	URL string

	Dependencies []string
	Conflicts    []string

	Tap     TapConfig
	Folder  string
	Install string
	Test    string
}

// Name return the name of the action.
func (*Action) Name() string {
	return "homebrew"
}

// GetInitConfig return what a good starting config would be.
func (*Action) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"homepage": "https://example.com",
		"url":      "https://github.com/owner/repo/releases/download/{{ .Version }}/darwin_amd64.zip",
		"tap": map[string]string{
			"url": "git@github.com:owner/homebrew-brew.git",
		},
	}
}

// Weight return in what order this action should be handled.
func (*Action) Weight() int {
	return 100
}

// Execute handles the action.
func (c *Action) Execute(cfg config.LetItGoConfig) error {
	setDefaults(c)
	templateProps(c, &cfg)
	hash, err := utils.BuildURLHash("sha256", c.URL)
	if err != nil {
		return err
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
	return utils.PublishFile(c.Tap.URL, filename, content, message)
}
