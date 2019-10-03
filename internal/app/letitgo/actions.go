package letitgo

import (
	"sort"

	"github.com/sirupsen/logrus"
)

var actions = Actions{}

type Actions []Action

func (s Actions) Len() int      { return len(s) }
func (s Actions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type Action interface {
	Weight() int
	Execute(Config) error
}

func RegisterAction(a Action) {
	actions = append(actions, a)
}

func RunAll(cfg Config) error {
	sort.Sort(ByWeight{actions})
	for _, a := range actions {
		logrus.Tracef("Running %T", a)
		if err := a.Execute(cfg); err != nil {
			return err
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
