package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/pkg/errors"
)

// New returns an action for archive
func New() action.Action {
	return &archive{}
}

type archive struct {
	Source string
	Target string
	Extras []string
	Method string
}

func (*archive) Name() string {
	return "archive"
}

func (*archive) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"source": "./build/bin/*",
		"target": "./build/pkg/",
		"extras": []string{"LICENSE.md"},
		"method": "zip",
	}
}

func (*archive) Weight() int {
	return 12
}

func (c *archive) Execute(cfg config.LetItGoConfig) error {
	if c.Method == "" {
		c.Method = "zip"
	}

	if _, err := os.Stat(c.Source); err != nil {
		return err
	}

	directories, err := filepath.Glob(c.Source)
	if err != nil {
		return fmt.Errorf("failed to resolve directories - %s", err.Error())
	}

	if err := os.MkdirAll(c.Target, 0755); err != nil {
		return fmt.Errorf("failed to create target directory - %s", err.Error())
	}

	for _, dir := range directories {
		target := path.Join(c.Target, filepath.Base(dir)+".zip")

		switch c.Method {
		case "zip":
			if err := zipCreate(dir, target, c.Extras); err != nil {
				return errors.Wrap(err, "zipCreate failed")
			}
		default:
			return errors.New("Unknown method")
		}
	}

	return nil
}

func zipCreate(source, target string, extras []string) error {
	for _, extra := range extras {
		if err := copyFile(
			path.Join(".", extra),
			path.Join(source, extra),
		); err != nil {
			return errors.Wrapf(err, "Failed to copy %s", extra)
		}
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer utils.DeferCheck(zipfile.Close)

	archive := zip.NewWriter(zipfile)
	defer utils.DeferCheck(archive.Close)

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if !info.IsDir() {
		return fmt.Errorf("archive needs source to be a directory")
	}
	baseDir = filepath.Base(source)

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if source == path {
			return err
		}
		return zipAddFile(archive, baseDir, info, path, source)
	})
}

func zipAddFile(archive *zip.Writer, baseDir string, info os.FileInfo, path, source string) error {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	if baseDir != "" {
		header.Name = filepath.Join(strings.TrimPrefix(path, source))[1:]
	}

	if info.IsDir() {
		header.Name += "/"
	} else {
		header.Method = zip.Deflate
	}

	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer utils.DeferCheck(file.Close)
	_, err = io.Copy(writer, file)

	return err
}

func copyFile(source, target string) error {
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	from, err := os.Open(source)
	if err != nil {
		return err
	}
	defer utils.DeferCheck(from.Close)

	to, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, info.Mode())
	if err != nil {
		return err
	}
	defer utils.DeferCheck(to.Close)

	if _, err = io.Copy(to, from); err != nil {
		return err
	}

	return nil
}
