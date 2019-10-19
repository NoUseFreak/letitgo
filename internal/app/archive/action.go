package archive

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/NoUseFreak/letitgo/internal/app/config"
)

// Action can package a directory into an archive.
type Action struct {
	Source string
	Target string
	Extras []string
	Method string
}

// Name return the name of the action.
func (*Action) Name() string {
	return "archive"
}

// GetInitConfig return what a good starting config would be.
func (*Action) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"source": "./build/bin/*",
		"target": "./build/pkg/",
		"extras": []string{"LICENSE.md"},
		"method": "zip",
	}
}

// Weight return in what order this action should be handled.
func (*Action) Weight() int {
	return 12
}

// Execute handles the action.
func (c *Action) Execute(cfg config.LetItGoConfig) error {
	if c.Method == "" {
		c.Method = "zip"
	}

	directories, err := filepath.Glob(c.Source)
	if err != nil {
		return fmt.Errorf("Failed to resolve directories - %s", err.Error())
	}

	if err := os.MkdirAll(c.Target, 0755); err != nil {
		return fmt.Errorf("Failed to create target directory - %s", err.Error())
	}

	for _, dir := range directories {
		target := path.Join(c.Target, filepath.Base(dir)+".zip")

		switch c.Method {
		case "zip":
			if err := zipCreate(dir, target, c.Extras); err != nil {
				return err
			}
		default:
			return errors.New("Unknown method")
		}
	}

	return nil
}

func zipCreate(source, target string, extras []string) error {
	for _, extra := range extras {
		copyFile(
			path.Join(".", extra),
			path.Join(source, extra),
		)
	}

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if !info.IsDir() {
		return fmt.Errorf("archive needs source to be a directory")
	}
	baseDir = filepath.Base(source)

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if source == path {
			return err
		}
		return zipAddFile(archive, baseDir, info, path, source)
	})

	return err
}

func zipAddFile(archive *zip.Writer, baseDir string, info os.FileInfo, path, source string) error {
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	if baseDir != "" {
		header.Name = filepath.Join(strings.TrimPrefix(path, source))
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
	defer file.Close()
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
	defer from.Close()

	to, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, info.Mode())
	if err != nil {
		return err
	}
	defer to.Close()

	if _, err = io.Copy(to, from); err != nil {
		return err
	}

	return nil
}
