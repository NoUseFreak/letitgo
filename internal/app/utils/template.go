package utils

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
)

func Template(tplStr string, data interface{}) (string, error) {
	tpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(tplStr)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}
