package changelog

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func buildReleaseBlocks(repo *git.Repository, ignore []string) (*[]releaseBlock, error) {
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("Failed to find HEAD - %s", err.Error())
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
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch log - %s", err.Error())
	}

	tree := []releaseBlock{{
		Tag:     "unreleased",
		Commits: []commit{},
	}}
	last := 0

	cIter.ForEach(func(c *object.Commit) error {
		if tag, ok := tags[c.Hash.String()]; ok {
			last++
			tree = append(tree, releaseBlock{
				Tag:     tag.Name().Short(),
				Date:    &c.Committer.When,
				Ref:     tag,
				Commits: []commit{},
			})
		}
		for _, ig := range ignore {
			if ig == c.Message {
				return nil
			}
		}
		for _, pc := range tree[last].Commits {
			if pc.Message == c.Message {
				return nil
			}
		}
		tree[last].Commits = append(tree[last].Commits, newCommit(c))
		return nil
	})

	return &tree, nil
}

func newCommit(c *object.Commit) commit {
	commit := commit{
		Message:      c.Message,
		MessageShort: strings.Split(utils.NormalizeNewlines(c.Message), "\n")[0],
		Hash:         c.Hash.String(),
	}

	return commit
}

type releaseBlock struct {
	Tag     string
	Date    *time.Time
	Ref     *plumbing.Reference
	Commits []commit
}

type commit struct {
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
