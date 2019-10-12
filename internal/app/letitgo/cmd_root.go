package letitgo

import (
	"fmt"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/config"
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
}

// Execute runs the cli application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	version := getVersion(args)
	logrus.Debugf("Going to work with version '%s'", version)

	cfg := Config{}
	cfg.LetItGo = config.NewConfig(version)
	utils.ParseYamlFile(".release.yml", &cfg)
	if err := RunAll(cfg); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func getVersion(args []string) string {
	if len(args) == 1 {
		return args[0]
	}

	v, err := utils.Run("describe", "--tags", "--abbrev", "0")
	if err != nil {
		logrus.Error("Could not find an exact tag.")
		os.Exit(0)
	}

	return v
}
