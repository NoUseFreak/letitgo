package letitgo

import (
	"fmt"
	"io/ioutil"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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
		ioutil.WriteFile(cfgFile, []byte(content), 0644)
		os.Exit(0)
	}

	ui.Phase("Printing config")
	fmt.Println("\n\n---\n" + content)
}

func getData() (string, error) {
	name := ""
	survey.AskOne(&survey.Input{
		Message: "Project name",
	}, &name)

	description := ""
	survey.AskOne(&survey.Input{
		Message: "Project description",
	}, &description)

	ligcfg := map[string]interface{}{
		"letitgo": map[string]string{
			"name":        name,
			"description": description,
		},
	}

	enabled := Actions{}
	for _, a := range actions {
		if askEnableAction(a.Name()) {
			enabled = append(enabled, a)
		}
	}

	out, err := yaml.Marshal(ligcfg)
	if err != nil {
		return "", err
	}
	content := string(out)

	for _, a := range enabled {
		d := map[string]interface{}{
			a.Name(): []interface{}{a.GetInitConfig()},
		}
		aout, err := yaml.Marshal(d)
		if err != nil {
			return "", err
		}
		content += "\n" + string(aout)
	}

	return content, nil
}

func askEnableAction(name string) bool {
	response := false
	survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Enable action %s", name),
		Default: true,
	}, &response)

	return response
}
