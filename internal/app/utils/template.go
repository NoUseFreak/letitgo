package utils

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/NoUseFreak/letitgo/internal/app/config"
)

// Template is a wrapper to standardize templating.
func Template(tplStr string, cfg config.LetItGoConfig, data ...interface{}) (string, error) {
	tpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(tplStr)
	if err != nil {
		return "", err
	}

	var tplVars map[string]interface{}
	inrec, _ := json.Marshal(cfg)
	json.Unmarshal(inrec, &tplVars)
	for _, d := range data {
		inrec, _ = json.Marshal(d)
		json.Unmarshal(inrec, &tplVars)
	}
	tplVars["Version"] = cfg.Version()

	var out bytes.Buffer
	if err := tpl.Execute(&out, tplVars); err != nil {
		return "", err
	}

	return out.String(), nil
}
