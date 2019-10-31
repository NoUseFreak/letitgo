package gofish

import (
	"fmt"
	"strings"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/NoUseFreak/letitgo/internal/app/utils/git"
	"github.com/pkg/errors"
)

const (
	gofishOrganization = "fishworks"
	gofishRepository   = "fish-food"
)

// New returns an action for githubrelease
func New() action.Action {
	return &gofish{}
}

type gofish struct {
	GithubUsername string            `yaml:"githubusername"`
	Homepage       string            `yaml:"homepage"`
	Artifacts      []*gofishArtifact `yaml:"artifacts"`
	CreatePR       bool
}

type gofishArtifact struct {
	Os     string `yaml:"os"`
	Arch   string `yaml:"arch"`
	URL    string `yaml:"url"`
	Sha256 string `yaml:"-"`
}

func (*gofish) Name() string {
	return "gofish"
}

func (*gofish) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"artifacts": []map[string]string{
			{
				"os":   "darwin",
				"arch": "amd64",
				"url":  "https://github.com/username/project/releases/download/{{ .Version }}/darwin_amd64.zip",
			},
			{
				"os":   "linux",
				"arch": "amd64",
				"url":  "https://github.com/username/project/releases/download/{{ .Version }}/linux_amd64.zip",
			},
			{
				"os":   "windows",
				"arch": "amd64",
				"url":  "https://github.com/username/project/releases/download/{{ .Version }}/windows_amd64.zip",
			},
		},
	}
}

func (*gofish) Weight() int {
	return 100
}

func (c *gofish) Execute(cfg config.LetItGoConfig) error {

	c.templateProps(&cfg)

	client, err := git.NewClient(
		fmt.Sprintf(
			"git@github.com:%s/%s.git",
			c.GithubUsername,
			gofishRepository,
		),
		utils.DryRun.IsEnabled(),
	)
	if err != nil {
		return err
	}

	var branchname string
	if !client.Exists() {
		ui.Step("Fork not found, creating...")
		if err := client.CreateForkFrom(gofishOrganization, gofishRepository); err != nil {
			return err
		}
		ui.Step("Fork created")
		branchname = fmt.Sprintf("add-%s", strings.ToLower(cfg.Name))
	} else {
		branchname = fmt.Sprintf("update-%s", cfg.Name)
	}

	ui.Step("Building artifact hashes")
	for _, a := range c.Artifacts {
		s, _ := utils.BuildURLHash("sha256", a.URL)
		a.Sha256 = s
	}

	path := fmt.Sprintf("Food/%s.lua", cfg.Name)
	message := fmt.Sprintf("%s %s", cfg.Name, cfg.Version())
	content, err := utils.Template(luaTpl, cfg, c)
	if err != nil {
		return errors.Wrap(err, "Failed to template lua file")
	}

	ui.Step("Publishing lua script")
	err = client.PublishFile(path, content, message, &branchname)
	if err != nil {
		return errors.Wrapf(err, "Failed to publish gofish file")
	}

	if !c.CreatePR {
		url := fmt.Sprintf(
			"https://github.com/%s/%s/compare/master...%s:%s?expand=true",
			gofishOrganization,
			gofishRepository,
			c.GithubUsername,
			branchname,
		)
		ui.Step("PR creation not enabled. Create PR: %s", url)
		return nil
	}

	ui.Step("Creating PR")

	ui.Error("Not implemented.")

	return nil
}

func (c *gofish) templateProps(cfg *config.LetItGoConfig) {
	for _, a := range c.Artifacts {
		utils.TemplateProperty(&a.URL, c, cfg)
	}
}
