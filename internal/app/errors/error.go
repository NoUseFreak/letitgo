package errors

import "fmt"

type SkipError struct {
	Part   string
	Reason string
}

func (e *SkipError) Error() string {
	return fmt.Sprintf("Skipping - %s (%s)", e.Part, e.Reason)
}
