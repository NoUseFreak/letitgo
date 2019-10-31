package git

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

// GetRemote looks at the given directory and returns the first git-remote name.
func GetRemote(dir string) (string, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return "", errors.Wrap(err, "Failed opening git repo")
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", errors.Wrap(err, "Failed reading remotes")
	}

	if len(remotes) > 1 {
		return "", errors.New("Can't resolve remote")
	}

	return remotes[0].Config().URLs[0], nil
}
