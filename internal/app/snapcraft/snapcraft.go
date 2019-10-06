package snapcraft

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/sirupsen/logrus"
)

func Execute(c Config) error {
	c.Name = utils.DefaultString(c.Name, c.LetItGo.Name)

	if _, err := exec.LookPath("snapcraft"); err != nil {
		return errors.New("snapcraft binary not installed")
	}

	dir := path.Join("/tmp", ".letitgo-"+c.Name)
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

	content, err := utils.Template(snapcraftTpl, &c)
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

	name := fmt.Sprintf("%s_%s_%s.snap", c.Name, c.Version(), c.Architecture)
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
