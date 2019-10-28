package letitgo

import (
	"fmt"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "letitgo <version>",
	Short: "LetItGo Release helper",
	Long:  `LetItGo release helper`,
	Args:  cobra.RangeArgs(0, 1),
	Run:   runRoot,
}

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	rootCmd.PersistentFlags().BoolP("dry-run", "d", false, "Enable dry-run")
	rootCmd.PersistentFlags().StringP("config", "c", ".release.yml", "Config file to use")
	rootCmd.PersistentFlags().String("loglevel", "info", "Log level")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if dryRun, _ := cmd.Flags().GetBool("dry-run"); dryRun {
			utils.DryRun.Enable()
		}

		lvlString, _ := cmd.PersistentFlags().GetString("loglevel")
		lvl, err := logrus.ParseLevel(lvlString)
		if err == nil {
			logrus.SetLevel(lvl)
		}
		return nil
	}
}

// Execute runs the cli application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	releaseCmd.Run(cmd, args)
}

func getVersion(args []string) string {
	if len(args) == 1 {
		return args[0]
	}

	v, err := utils.Run("describe", "--tags", "--abbrev")
	if err != nil {
		ui.Error("Could not find an exact tag.")
		os.Exit(0)
	}

	return v
}
