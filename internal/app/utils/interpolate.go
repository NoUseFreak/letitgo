package utils

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
)

func TemplateProperty(v *string, ctx interface{}) {
	if tpl, err := template.New("value").Funcs(sprig.TxtFuncMap()).Parse(*v); err == nil {
		var out bytes.Buffer
		if err := tpl.Execute(&out, ctx); err == nil {
			*v = out.String()
		}
	}
}
