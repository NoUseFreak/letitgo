package utils

import "github.com/NoUseFreak/letitgo/internal/app/ui"

func DeferCheck(fs ...func() error) {
	for i := len(fs) - 1; i >= 0; i-- {
		if err := fs[i](); err != nil {
			ui.Error("Received error:", err)
		}
	}
}
