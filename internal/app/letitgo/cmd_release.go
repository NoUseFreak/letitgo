package letitgo

import (
	"os"
	"sort"

	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release",
	Long:  `release`,
	Run:   executeRelease,
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}

func executeRelease(cmd *cobra.Command, args []string) {
	ui.Title("LetItGo")
	cfgWrapper := struct{ LetItGo Config }{LetItGo: Config{}}
	utils.ParseYamlFile(".release.yml", &cfgWrapper)
	cfg := cfgWrapper.LetItGo

	var workload Actions
	for _, a := range cfg.Actions {
		actionType := a["type"].(string)
		if action := getActions()[actionType]; action != nil {
			mapstructure.Decode(a, action)
			mapstructure.Decode(cfg, action)

			workload = append(workload, action)
		}
	}

	sort.Sort(ByWeight{workload})

	version := getVersion(args)

	ligConfig := config.NewConfig(version)
	ligConfig.Name = cfg.Name
	ligConfig.Description = cfg.Description

	if len(workload) == 0 {
		color.Yellow("No supported actions found")
		os.Exit(0)
	}

	for _, action := range workload {
		ui.Phase("Releasing %s", action.Name())
		if err := action.Execute(ligConfig); err != nil {
			switch er := err.(type) {
			case *e.SkipError:
				color.Yellow("  " + er.Error())
			default:
				ui.Error(err.Error())
				os.Exit(1)
			}
		}
	}
}
