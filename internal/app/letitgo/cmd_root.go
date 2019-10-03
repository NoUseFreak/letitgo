package letitgo

import (
	"fmt"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "letitgo <version>",
	Short: "LetItGo Release helper",
	Long:  `LetItGo release helper`,
	Args:  cobra.ExactArgs(1),
	Run:   runRoot,
}

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

// Execute runs the cli application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	cfg := Config{}
	utils.ParseYamlFile(".release.yml", &cfg)
	cfg.Version = args[0]
	if err := RunAll(cfg); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
