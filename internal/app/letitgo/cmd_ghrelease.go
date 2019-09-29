package letitgo

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ghreleaseCmd = &cobra.Command{
	Use:   "ghrelease <version>",
	Short: "Release to github",
	Long:  `Release to github`,
	Args:  cobra.ExactArgs(1),
	Run:   runGhRelease,
}

func init() {
	rootCmd.AddCommand(ghreleaseCmd)
}

func runGhRelease(cmd *cobra.Command, args []string) {
	cfg := NewConfig("./.release.yml")

	if cfg.GhRelease == nil {
		logrus.Error("No ghrelease config found")
		os.Exit(1)
	}

	for _, task := range cfg.GhRelease {
		task.Version = args[0]
		task.Defaults()
		if err := task.Execute(); err != nil {
			logrus.Error(err)
		}
	}
}
