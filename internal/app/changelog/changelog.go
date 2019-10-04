package changelog

import (
	"bytes"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Execute(c Config) error {
	if c.File == "" {
		c.File = "CHANGELOG.md"
	}
	if c.Message == "" {
		c.Message = "Update changelog"
	}

	tree, err := buildReleaseBlocks()
	if err != nil {
		return err
	}

	vars := struct {
		Blocks []ReleaseBlock
	}{
		Blocks: *tree,
	}
	out, err := templateChangelog(vars)
	if err != nil {
		return err
	}

	repo, err := utils.GetRemote(".")
	if err != nil {
		return err
	}

	logrus.Infof("Publishing %s", c.File)
	return utils.PublishFile(repo, c.File, out, c.Message)
}

func buildReleaseBlocks() (*[]ReleaseBlock, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	tags := map[string]*plumbing.Reference{}
	tagrefs, _ := repo.Tags()
	tagrefs.ForEach(func(t *plumbing.Reference) error {
		tags[t.Hash().String()] = t
		return nil
	})

	cIter, err := repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})

	tree := []ReleaseBlock{
		{
			Tag:     "unreleased",
			Commits: []*object.Commit{},
		},
	}
	last := 0

	err = cIter.ForEach(func(c *object.Commit) error {
		if tag, ok := tags[c.Hash.String()]; ok {
			last++
			tree = append(tree, ReleaseBlock{
				Tag:     tag.Name().Short(),
				Date:    &c.Committer.When,
				Ref:     tag,
				Commits: []*object.Commit{},
			})
		}
		tree[last].Commits = append(tree[last].Commits, c)
		return nil
	})

	return &tree, err
}

type ReleaseBlock struct {
	Tag     string
	Date    *time.Time
	Ref     *plumbing.Reference
	Commits []*object.Commit
}

func templateChangelog(vars interface{}) (string, error) {
	tpl := template.Must(template.New("main").Parse(changelogTpl))
	tpl.Funcs(sprig.TxtFuncMap())
	template.Must(tpl.New("block").Parse(blockTpl))

	var out bytes.Buffer
	if err := tpl.Execute(&out, vars); err != nil {
		return "", err
	}

	return out.String(), nil
}
