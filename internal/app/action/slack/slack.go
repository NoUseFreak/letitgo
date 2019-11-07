package slack

import (
	"errors"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/utils"

	gitclient "github.com/NoUseFreak/letitgo/internal/app/utils/git"
	slackclient "github.com/nlopes/slack"
)

// New returns an action for archive
func New() action.Action {
	return &slack{}
}

type slack struct {
	Channel string

	Author  string
	Message string
}

func (*slack) Name() string {
	return "slack"
}

func (*slack) GetInitConfig() map[string]interface{} {
	return map[string]interface{}{
		"channel": "#released",
		"message": "Project released",
	}
}

func (*slack) Weight() int {
	return 99999
}

func (c *slack) Execute(cfg config.LetItGoConfig) error {
	token := os.Getenv("SLACK_TOKEN")
	if token == "" {
		return errors.New("make sure to set SLACK_TOKEN")
	}

	c.setDefaults()
	utils.TemplateProperty(&c.Message, c, &cfg)

	return c.sendMessage(token)
}

func (c *slack) setDefaults() {
	if c.Author == "" {
		c.Author = "LetItGo"
	}
	if c.Message == "" {
		c.Message = "Project release completed"
	}
	if c.Channel == "" {
		c.Channel = "#released"
	}
}

func (c *slack) sendMessage(token string) error {
	repo, err := gitclient.GetRemote(".")
	if err != nil {
		return err
	}

	api := slackclient.New(token)
	attachment := slackclient.Attachment{
		Color:      "good",
		AuthorName: c.Author,
		AuthorLink: "https://github.com/NoUseFreak/letitgo",
		Text:       c.Message,
		Fields: []slackclient.AttachmentField{
			slackclient.AttachmentField{
				Title: "repo",
				Value: repo,
			},
		},
	}

	_, _, err = api.PostMessage(
		c.Channel,
		slackclient.MsgOptionUsername(c.Author),
		slackclient.MsgOptionIconEmoji(":rocket:"),
		slackclient.MsgOptionAttachments(attachment),
	)

	return err
}
