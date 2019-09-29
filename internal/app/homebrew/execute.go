package homebrew

import (
	"fmt"
	"path"
)

func (h *Homebrew) Execute() error {
	h.templateInput()
	hash, err := h.buildURLHash(h.URL)
	if err != nil {
		return err
	}
	h.Hash = hash

	content, err := h.template()
	if err != nil {
		return err
	}
	filename := path.Join(h.Folder, fmt.Sprintf("%s.rb", h.Name))

	if err := h.saveToGit(filename, content); err != nil {
		return err
	}

	return nil
}
