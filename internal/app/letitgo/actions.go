package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/changelog"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/githubrelease"
	"github.com/NoUseFreak/letitgo/internal/app/helm"
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
	"github.com/NoUseFreak/letitgo/internal/app/snapcraft"
)

func init() {
	registerAction(new(changelog.Action))
	registerAction(new(githubrelease.Action))
	registerAction(new(helm.Action))
	registerAction(new(homebrew.Action))
	registerAction(new(snapcraft.Action))
}

var letItGoActions = map[string]Action{}

// Actions is a slice of Action structs.
type Actions []Action

// Action is an action LetItGo can handle.
type Action interface {
	Name() string
	GetInitConfig() map[string]interface{}
	// GetDefaults() Config
	Weight() int
	Execute(config.LetItGoConfig) error
}

// Len returns the length of the slice.
func (s Actions) Len() int { return len(s) }

// Swap swaps two items from possition.
func (s Actions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func registerAction(action Action) {
	letItGoActions[action.Name()] = action
}

func getActions() map[string]Action {
	return letItGoActions
}

// ByWeight allows to sort Actions by Weight.
type ByWeight struct {
	Actions
}

// Less Compares weights of each action.
func (s ByWeight) Less(i, j int) bool {
	return s.Actions[i].Weight() < s.Actions[j].Weight()
}
