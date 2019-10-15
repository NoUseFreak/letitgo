package letitgo

import (
	"sort"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/fatih/color"
)

var actions = Actions{}

type Actions []Action

func (s Actions) Len() int      { return len(s) }
func (s Actions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type Action interface {
	Name() string
	GetInitConfig() map[string]interface{}
	GetDefaults() Config
	Weight() int
	Execute(Config) error
}

func registerAction(a Action) {
	actions = append(actions, a)
}

func RunAll(cfg Config) error {
	sort.Sort(ByWeight{actions})
	for _, a := range actions {
		ui.Trace("Running %T", a)
		if err := a.Execute(cfg); err != nil {
			switch er := err.(type) {
			case *e.SkipError:
				color.Yellow("  " + er.Error())
			default:
				return err
			}
		}
	}

	return nil
}

type ByWeight struct {
	Actions
}

func (s ByWeight) Less(i, j int) bool {
	return s.Actions[i].Weight() < s.Actions[j].Weight()
}
