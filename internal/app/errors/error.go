package errors

import "fmt"

type SkipError struct {
	Reason string
}

func (e *SkipError) Error() string {
	return fmt.Sprintf("Skipping because of %s", e.Reason)
}
