package utils

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

// Run executes a git command and returns the result.
func Run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	bts, err := cmd.CombinedOutput()
	logrus.WithField("output", string(bts)).Trace("git result")
	if err != nil {
		return "", errors.New(string(bts))
	}

	output := strings.Replace(strings.Split(string(bts), "\n")[0], "'", "", -1)
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
