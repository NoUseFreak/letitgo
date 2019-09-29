package ghrelease

func (task *GhRelease) Defaults() {
	if task.Title == "" {
		task.Title = "{{ .Version }}"
	}
	if task.Description == "" {
		task.Description = "{{ .Version }}"
	}
}
