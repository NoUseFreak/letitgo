package changelog

import (
	"fmt"

	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"gopkg.in/src-d/go-git.v4"
)

// Action calculates a changelog for the current project.
type Action struct {
	File    string
	Message string
}

// Name return the name of the action.
func (a *Action) Name() string {
	return "changelog"
}

// GetInitConfig return what a good starting config would be.
func (a *Action) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"file": "CHANGELOG.md",
	}
}

// Weight return in what order this action should be handled.
func (a *Action) Weight() int {
	return 5
}

// Execute handles the action.
func (a *Action) Execute(cfg config.LetItGoConfig) error {
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
	return utils.PublishFile(repo, a.File, out, a.Message)
}
