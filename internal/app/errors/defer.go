package errors

import "github.com/NoUseFreak/letitgo/internal/app/ui"

// DeferCheck wraps a deferred function and prints the error if any.
func DeferCheck(fs ...func() error) {
	for i := len(fs) - 1; i >= 0; i-- {
		if err := fs[i](); err != nil {
			ui.Error("Received error: %s", err)
		}
	}
}
