package docker

import (
	"fmt"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	dckr "github.com/fsouza/go-dockerclient"
)

// New returns an action for docker
func New() action.Action {
	return &docker{}
}

type docker struct {
	Dockerfile string
	Images     []string
	NoPush     bool
	Auth       dockerAuth `yaml:"auth"`
}

type dockerAuth struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty" env:"DOCKER_AUTH_PASSWORD"`
}

func (*docker) Name() string {
	return "docker"
}

func (*docker) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"dockerfile": "Dockerfile",
		"images": []string{
			"user/app:{{ .Version }}",
			"user/app:{{ .Version.Major }}",
			"user/app:{{ .Version.Major }}.{{ .Version.Minor }}",
		},
	}
}

func (*docker) Weight() int {
	return 20
}

func (c *docker) Execute(cfg config.LetItGoConfig) error {
	var imageNames []string
	for _, template := range c.Images {
		name, err := utils.Template(template, cfg, c)
		if err != nil {
			return fmt.Errorf("Failed to build image names - %s", err.Error())
		}
		imageNames = append(imageNames, name)
	}

	client, err := dckr.NewClientFromEnv()
	if err != nil {
		return err
	}

	ui.Step("Building image")
	if err := c.buildImage(client, imageNames[0]); err != nil {
		return err
	}

	ui.Step("Tagging image %s", imageNames[0])
	if err := c.tagImages(client, imageNames[0], imageNames[1:]); err != nil {
		return err
	}

	if utils.DryRun.IsEnabled() {
		return &e.SkipError{
			Reason: "dryrun",
			Part:   "push images",
		}
	}

	if c.NoPush {
		return nil
	}

	ui.Step("Pushing images")
	c.pushImages(client, imageNames)
	return nil
}

func (c *docker) buildImage(client *dckr.Client, imageName string) error {
	wd, _ := os.Getwd()
	r, w := newProgressWriter()
	defer r.Close()
	defer w.Close()

	return client.BuildImage(dckr.BuildImageOptions{
		ContextDir:   wd,
		Dockerfile:   c.Dockerfile,
		Name:         imageName,
		OutputStream: w,
	})
}

func (c *docker) tagImages(client *dckr.Client, baseName string, imageNames []string) error {
	for _, name := range imageNames {
		ui.Step("Tagging image %s", name)
		img := parseImageName(name)
		err := client.TagImage(baseName, dckr.TagImageOptions{
			Repo: img.Repo,
			Tag:  img.Tag,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *docker) pushImages(client *dckr.Client, imageNames []string) error {
	var errs []error
	r, w := newProgressWriter()
	defer r.Close()
	defer w.Close()
	for _, name := range imageNames {
		img := parseImageName(name)
		err := client.PushImage(dckr.PushImageOptions{
			Name:         img.Repo,
			Tag:          img.Tag,
			OutputStream: w,
		}, dckr.AuthConfiguration{
			Username: c.Auth.Username,
			Password: c.Auth.Password,
		})
		if err != nil {
			ui.Warn("Failed to push image - %v", err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
