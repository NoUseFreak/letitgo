package helm

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

// Action can package and publish helm charts.
type Action struct {
	Charts     []string
	BuildDir   string
	Repository string
}

// Name return the name of the action.
func (*Action) Name() string {
	return "helm"
}

// GetInitConfig return what a good starting config would be.
func (*Action) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"charts":     []string{"./charts/my-chart"},
		"builddir":   "./build/helm",
		"repository": "https://chartmuseum.example.com",
	}
}

// Weight return in what order this action should be handled.
func (*Action) Weight() int {
	return 120
}

// Execute handles the action.
func (c *Action) Execute(cfg config.LetItGoConfig) error {
	if c.BuildDir == "" {
		c.BuildDir = "./build/helm"
	}

	if _, err := exec.LookPath("helm"); err != nil {
		return errors.New("helm binary not installed")
	}

	helm := utils.Helm{}
	if len(c.Charts) != 0 {
		for _, chart := range c.Charts {
			if err := helm.Package(chart, c.BuildDir, cfg.Version()); err != nil {
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
