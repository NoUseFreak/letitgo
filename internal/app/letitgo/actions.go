package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/action/archive"
	"github.com/NoUseFreak/letitgo/internal/app/action/changelog"
	"github.com/NoUseFreak/letitgo/internal/app/action/docker"
	"github.com/NoUseFreak/letitgo/internal/app/action/githubrelease"
	"github.com/NoUseFreak/letitgo/internal/app/action/gofish"
	"github.com/NoUseFreak/letitgo/internal/app/action/helm"
	"github.com/NoUseFreak/letitgo/internal/app/action/homebrew"
	"github.com/NoUseFreak/letitgo/internal/app/action/slack"
	"github.com/NoUseFreak/letitgo/internal/app/action/snapcraft"
)

var _letItGoActions = map[string]action.Action{}

func init() {
	registerAction(archive.New())
	registerAction(changelog.New())
	registerAction(githubrelease.New())
	registerAction(gofish.New())
	registerAction(helm.New())
	registerAction(homebrew.New())
	registerAction(slack.New())
	registerAction(snapcraft.New())
	registerAction(docker.New())
}

func registerAction(action action.Action) {
	_letItGoActions[action.Name()] = action
}

func getActions() map[string]action.Action {
	return _letItGoActions
}
