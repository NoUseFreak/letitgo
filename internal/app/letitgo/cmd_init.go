package letitgo

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	survey "github.com/AlecAivazis/survey/v2"
	yaml "gopkg.in/yaml.v2"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init LetItGo config file",
	Long:  `Init`,
	Run:   executeInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func executeInit(cmd *cobra.Command, args []string) {
	ui.Title("LetItGo")

	content, err := getData()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	cfgFile := "./.release.yml"
	_, err = os.Stat(cfgFile)
	if os.IsNotExist(err) {
		ui.Phase("Creating config")
		if err := ioutil.WriteFile(cfgFile, []byte(content), 0644); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

	ui.Phase("Printing config")
	fmt.Println("\n\n---\n" + content)
}

func getData() (string, error) {
	name := ""
	_ = survey.AskOne(&survey.Input{
		Message: "Project name",
	}, &name)

	description := ""
	_ = survey.AskOne(&survey.Input{
		Message: "Project description",
	}, &description)

	enabled := []yaml.MapSlice{}
	for _, a := range getActions() {
		if askEnableAction(a.Name()) {
			action := yaml.MapSlice{{
				Key:   "type",
				Value: a.Name(),
			}}
			for k, v := range a.GetInitConfig() {
				action = append(action, yaml.MapItem{
					Key:   k,
					Value: v,
				})
			}
			enabled = append(enabled, action)
		}
	}

	ligcfg := map[string]interface{}{
		"letitgo": yaml.MapSlice{
			{Key: "name", Value: name},
			{Key: "description", Value: description},
			{Key: "actions", Value: enabled},
		},
	}

	out, err := yaml.Marshal(ligcfg)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func askEnableAction(name string) bool {
	response := false
	_ = survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Enable action %s", name),
		Default: true,
	}, &response)

	return response
}
