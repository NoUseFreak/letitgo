package errors

import "fmt"

// SkipError indicates an action was skipped.
type SkipError struct {
	Part   string
	Reason string
}

// Error return a string explaining why it was skipped.
func (e *SkipError) Error() string {
	return fmt.Sprintf("Skipping - %s (%s)", e.Part, e.Reason)
}
