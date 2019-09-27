package homebrew

type Homebrew struct {
	Name        string
	Description string
	Homepage    string

	URL     string
	Version string

	Dependencies []string
	Conflicts    []string

	Tap     HomebrewTap
	Folder  string
	Install string
	Test    string

	// Calculated
	Hash string
}

type HomebrewTap struct {
	URL         string
	AuthorName  string
	AuthorEmail string
}
