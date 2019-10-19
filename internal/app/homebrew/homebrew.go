package homebrew

import (
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

func setDefaults(c *Action) {
	if c.Folder == "" {
		c.Folder = "Formula"
	}
	if c.Install == "" {
		c.Install = "bin.install \"{{ .Name }}\""
	}
}

func templateProps(c *Action, cfg *config.LetItGoConfig) {
	utils.TemplateProperty(&c.URL, c, cfg)
	utils.TemplateProperty(&c.Install, c, cfg)
	utils.TemplateProperty(&c.Test, c, cfg)
}
