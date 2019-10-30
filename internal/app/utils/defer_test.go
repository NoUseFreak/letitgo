package utils

import (
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/pkg/errors"
)

func TestDeferCheck(t *testing.T) {
	fOk := func() error { return nil }
	fNok := func() error { return errors.New("Failed") }
	tests := []struct {
		name       string
		fs         []func() error
		wantOutput bool
	}{
		{
			name: "noerror",
			fs: []func() error{
				fOk,
			},
			wantOutput: false,
		},
		{
			name: "error",
			fs: []func() error{
				fNok,
			},
			wantOutput: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := capturer.CaptureOutput(func() {
				DeferCheck(tt.fs...)
			})

			if !tt.wantOutput && out != "" {
				t.Errorf("Did not expect output, got '%s'", out)
			}
			if tt.wantOutput && out == "" {
				t.Errorf("Expected output, got ''")
			}
		})
	}

}
