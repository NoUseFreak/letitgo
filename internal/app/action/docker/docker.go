package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/pkg/errors"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	dckrTypes "github.com/docker/docker/api/types"
	dckr "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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

	client, err := dckr.NewClientWithOpts(dckr.WithAPIVersionNegotiation(), dckr.FromEnv)
	if err != nil {
		panic(err)
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
	return c.pushImages(client, imageNames)
}

func (c *docker) buildImage(client *dckr.Client, imageName string) error {
	wd, _ := os.Getwd()

	resp, err := client.ImageBuild(context.Background(), getContext(wd), dckrTypes.ImageBuildOptions{
		SuppressOutput: false,
		PullParent:     true,
		Remove:         true,
		Dockerfile:     c.Dockerfile,
		Tags:           []string{imageName},
	})
	if err != nil {
		return errors.Wrapf(err, "Failed to build %s", c.Dockerfile)
	}
	defer utils.DeferCheck(resp.Body.Close)

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return err
}

func getContext(filePath string) io.Reader {
	ctx, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})
	return ctx
}

func (c *docker) tagImages(client *dckr.Client, baseName string, imageNames []string) error {
	for _, name := range imageNames {
		ui.Step("Tagging image %s", name)
		if err := client.ImageTag(context.Background(), baseName, name); err != nil {
			return err
		}
	}

	return nil
}

func (c *docker) pushImages(client *dckr.Client, imageNames []string) error {
	authConfig := dckrTypes.AuthConfig{
		Username: c.Auth.Username,
		Password: c.Auth.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	for _, name := range imageNames {
		out, err := client.ImagePush(context.Background(), name, dckrTypes.ImagePushOptions{
			RegistryAuth: authStr,
		})
		if err != nil {
			return err
		}
		defer utils.DeferCheck(out.Close)
		_, err = ioutil.ReadAll(out)
		if err != nil {
			return err
		}
	}

	return nil
}
