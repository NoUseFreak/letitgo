package action

// Actions is a slice of Actions.
type Actions []Action

// Len returns the length of the slice.
func (s Actions) Len() int { return len(s) }

// Swap swaps two items from possition.
func (s Actions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
