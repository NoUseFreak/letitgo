package changelog

import (
	"fmt"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"gopkg.in/src-d/go-git.v4"
)

// New returns an action for changelog
func New() action.Action {
	return &changelog{}
}

type changelog struct {
	File    string
	Message string
}

func (a *changelog) Name() string {
	return "changelog"
}

func (a *changelog) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"file": "CHANGELOG.md",
	}
}

func (a *changelog) Weight() int {
	return 5
}

func (a *changelog) Execute(cfg config.LetItGoConfig) error {
	if a.File == "" {
		a.File = "CHANGELOG.md"
	}
	if a.Message == "" {
		a.Message = "Update changelog\n[skip ci]"
	}

	r, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("Unable to find git repo - %s", err.Error())
	}

	if lastCommitIsChangelog(r, a.Message, a.File) {
		ui.Info("Skipping changelog")
		return nil
	}

	tree, err := buildReleaseBlocks(r, []string{a.Message})
	if err != nil {
		return fmt.Errorf("Unable to build release blocks - %s", err.Error())
	}

	vars := struct {
		Blocks []releaseBlock
	}{
		Blocks: *tree,
	}
	out, err := templateChangelog(vars)
	if err != nil {
		return fmt.Errorf("Unable to template changelog - %s", err.Error())
	}

	repo, err := utils.GetRemote(".")
	if err != nil {
		return fmt.Errorf("Unable to resolve remote - %s", err.Error())
	}

	ui.Trace(out)

	ui.Step("Publishing %s", a.File)

	utils.WriteFile(a.File, out)
	return utils.PublishFile(repo, a.File, out, a.Message)
}
