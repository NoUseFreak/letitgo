package helm

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

func Execute(c Config) error {
	ui.Phase("Helm")

	if c.BuildDir == "" {
		c.BuildDir = "./build/helm"
	}

	if _, err := exec.LookPath("helm"); err != nil {
		return errors.New("helm binary not installed")
	}

	helm := utils.Helm{}
	if len(c.Charts) != 0 {
		for _, chart := range c.Charts {
			if err := helm.Package(chart, c.BuildDir, c.BaseConfig.LetItGo.Version()); err != nil {
				return err
			}
		}
	}

	if c.Repository != "" {
		artifacts, err := filepath.Glob(c.BuildDir + "/*")
		if err != nil {
			return err
		}

		for _, artifact := range artifacts {
			if err := helm.Publish(artifact, c.Repository); err != nil {
				return err
			}
		}
	}

	return nil
}
