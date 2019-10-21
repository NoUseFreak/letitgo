package snapcraft

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/sirupsen/logrus"
)

// New returns an action for snapcraft
func New() action.Action {
	return &snapcraft{}
}

type snapcraft struct {
	Assets       []string
	Architecture string
}

func (*snapcraft) Name() string {
	return "snapcraft"
}

func (*snapcraft) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"assets": []string{
			"./build/bin/linux_amd64/letitgo",
		},
		"architecture": "amd64",
	}
}

func (*snapcraft) Weight() int {
	return 110
}

func (c *snapcraft) Execute(cfg config.LetItGoConfig) error {
	if _, err := exec.LookPath("snapcraft"); err != nil {
		return errors.New("snapcraft binary not installed")
	}

	dir := path.Join("/tmp", ".letitgo-"+cfg.Name)
	metaDir := path.Join(dir, "meta")

	if err := os.MkdirAll(metaDir, 0777); err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	for _, assetGlob := range c.Assets {
		files, _ := filepath.Glob(assetGlob)
		for _, path := range files {
			destPath := filepath.Join(dir, filepath.Base(path))
			os.Link(path, destPath)
		}
	}

	content, err := utils.Template(snapcraftTpl, cfg, c)
	if err != nil {
		return err
	}
	logrus.Trace(content)

	if err = ioutil.WriteFile(path.Join(dir, "meta", "snap.yaml"), []byte(content), 0644); err != nil {
		return err
	}

	cmd := exec.Command("snapcraft", "pack", dir)
	cmd.Dir = dir
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin
	logrus.Debug("Building snap")
	if err = cmd.Run(); err != nil {
		return err
	}

	name := fmt.Sprintf("%s_%s_%s.snap", cfg.Name, cfg.Version(), c.Architecture)
	cmd = exec.Command("snapcraft", "push", name)
	cmd.Dir = dir
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
