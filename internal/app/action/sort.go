package action

// ByWeight allows to sort Actions by Weight.
type ByWeight struct {
	Actions
}

// Less Compares weights of each action.
func (s ByWeight) Less(i, j int) bool {
	return s.Actions[i].Weight() < s.Actions[j].Weight()
}
