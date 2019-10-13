package changelog

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Execute(c Config) error {
	if c.File == "" {
		c.File = "CHANGELOG.md"
	}
	if c.Message == "" {
		c.Message = "Update changelog\n[skip ci]"
	}

	r, err := git.PlainOpen(".")
	if err != nil {
		return err
	}

	if lastCommitIsChangelog(r, c.Message, c.File) {
		ui.Info("Skipping changelog")
		return nil
	}

	tree, err := buildReleaseBlocks(r, []string{c.Message})
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

	ui.Trace(out)

	ui.Step("Publishing %s", c.File)
	return utils.PublishFile(repo, c.File, out, c.Message)
}

func buildReleaseBlocks(repo *git.Repository, ignore []string) (*[]ReleaseBlock, error) {
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
			Commits: []Commit{},
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
				Commits: []Commit{},
			})
		}
		for _, ig := range ignore {
			if ig == c.Message {
				return nil
			}
		}
		tree[last].Commits = append(tree[last].Commits, newCommit(c))
		return nil
	})

	return &tree, err
}

func newCommit(c *object.Commit) Commit {
	commit := Commit{
		Message:      c.Message,
		MessageShort: strings.Split(utils.NormalizeNewlines(c.Message), "\n")[0],
		Hash:         c.Hash.String(),
	}

	return commit
}

type ReleaseBlock struct {
	Tag     string
	Date    *time.Time
	Ref     *plumbing.Reference
	Commits []Commit
}

type Commit struct {
	Message      string
	MessageShort string
	Hash         string
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

func lastCommitIsChangelog(r *git.Repository, msg, file string) bool {
	ref, err := r.Head()
	if err != nil {
		return false
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return false
	}

	if f, _ := commit.File(file); f == nil {
		return false
	}

	return commit.Message == msg
}
