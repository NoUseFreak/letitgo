package homebrew

func (h *Homebrew) Defaults() {
	if h.Folder == "" {
		h.Folder = "Formula"
	}
	if h.Install == "" {
		h.Install = "bin.install \"{{ .Name }}\""
	}
}
