package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"

	"github.com/NoUseFreak/letitgo/internal/app/ui"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
)

// Helm provices a wrapper for all helm interactions.
type Helm struct{}

// Package uses the helm binary to package a chart.
func (h *Helm) Package(chart, target, version string) error {
	ui.Step("Packaging %s@%s", path.Base(chart), version)

	if _, err := os.Stat(target); os.IsNotExist(err) {
		os.MkdirAll(target, 0755)
	}
	if DryRun {
		return &e.SkipError{
			Reason: "dryrun",
			Part:   fmt.Sprintf("Packaging %s", path.Base(chart)),
		}
	}
	cmd := exec.Command("helm", "package", "-d", target, "--version", version, chart)
	return cmd.Run()
}

// Publish pushes an packages chart to a given repository.
func (h *Helm) Publish(artifact, repository string) error {
	ui.Step("Publishing %s", path.Base(artifact))

	u, err := url.Parse(repository)
	if err != nil {
		return err
	}
	u.Path = "/api/charts"

	if DryRun {
		return &e.SkipError{
			Reason: "dryrun",
			Part:   fmt.Sprintf("Publishing %s", path.Base(artifact)),
		}
	}

	return sendPostRequest(u.String(), artifact, "application/x-tgz")
}

func sendPostRequest(url string, filename string, filetype string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	_, err = http.Post(url, filetype, file)
	if err != nil {
		return err
	}

	return nil
}
