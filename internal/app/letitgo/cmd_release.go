package letitgo

import (
	"os"
	"sort"

	"github.com/NoUseFreak/letitgo/internal/app/action"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	env "github.com/caarlos0/env/v6"
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
	cfgFile, _ := cmd.Flags().GetString("config")
	utils.ParseYamlFile(cfgFile, &cfgWrapper)
	cfg := cfgWrapper.LetItGo

	if err := doRelease(cfg, getVersion(args)); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func doRelease(cfg Config, version string) error {
	var workload action.Actions
	for _, a := range cfg.Actions {
		actionType := a["type"].(string)
		if action := getActions()[actionType]; action != nil {
			if err := mapstructure.Decode(a, action); err != nil {
				ui.Error(err.Error())
			}
			if err := mapstructure.Decode(cfg, action); err != nil {
				ui.Error(err.Error())
			}
			if err := env.Parse(action); err != nil {
				ui.Warn("Failed to parse environment - %s", err.Error())
			}

			workload = append(workload, action)
		}
	}

	sort.Sort(action.ByWeight{Actions: workload})

	ligConfig := config.NewConfig(version)
	ligConfig.Name = cfg.Name
	ligConfig.Description = cfg.Description

	if len(workload) == 0 {
		color.Yellow("No supported actions found")
		return nil
	}

	for _, action := range workload {
		ui.Phase("Releasing %s", action.Name())
		if err := action.Execute(ligConfig); err != nil {
			switch er := err.(type) {
			case *e.SkipError:
				color.Yellow("  " + er.Error())
			default:
				ui.Error(err.Error())
				return err
			}
		}
	}

	return nil
}
