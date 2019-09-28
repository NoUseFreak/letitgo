package letitgo

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var homebrewCmd = &cobra.Command{
	Use:   "homebrew <version>",
	Short: "Update homebrew",
	Long:  `Update homebrew`,
	Args:  cobra.ExactArgs(1),
	Run:   runHomebrew,
}

func init() {
	rootCmd.AddCommand(homebrewCmd)
}

func runHomebrew(cmd *cobra.Command, args []string) {
	cfg := NewConfig("./.release.yml")

	if cfg.Homebrew == nil {
		logrus.Error("No homebrew config found")
		os.Exit(1)
	}

	for _, h := range cfg.Homebrew {
		h.Version = args[0]
		h.Defaults()
		if err := h.Execute(); err != nil {
			logrus.Error(err)
		}
	}
}
