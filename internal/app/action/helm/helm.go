package helm

import (
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
)

// New returns an action for helm
func New() action.Action {
	return &helm{}
}

type helm struct {
	Charts     []string
	BuildDir   string
	Repository string
}

func (*helm) Name() string {
	return "helm"
}

func (*helm) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"charts":     []string{"./charts/my-chart"},
		"builddir":   "./build/helm",
		"repository": "https://chartmuseum.example.com",
	}
}

func (*helm) Weight() int {
	return 120
}

func (c *helm) Execute(cfg config.LetItGoConfig) error {
	if c.BuildDir == "" {
		c.BuildDir = "./build/helm"
	}

	if _, err := exec.LookPath("helm"); err != nil {
		return errors.New("helm binary not installed")
	}

	helm := utils.Helm{}
	if len(c.Charts) != 0 {
		for _, chart := range c.Charts {
			if err := helm.Package(chart, c.BuildDir, cfg.Version().String()); err != nil {
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
