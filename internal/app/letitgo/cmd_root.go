package letitgo

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "letitgo <version>",
	Short: "LetItGo Release helper",
	Long:  `LetItGo`,
}

// Execute runs the cli application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
