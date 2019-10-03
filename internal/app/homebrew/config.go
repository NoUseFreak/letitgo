package homebrew

type Config struct {
	Name        string
	Description string
	Homepage    string

	URL     string
	Version string

	Dependencies []string
	Conflicts    []string

	Tap     HomebrewTapConfig
	Folder  string
	Install string
	Test    string

	// internal
	Hash string
}

type HomebrewTapConfig struct {
	URL         string
	AuthorName  string
	AuthorEmail string
}
