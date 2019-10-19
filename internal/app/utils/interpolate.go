package utils

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/NoUseFreak/letitgo/internal/app/config"
)

// TemplateProperty takes a property as a template and interpolates any variable.
func TemplateProperty(v *string, ctx interface{}, cfg *config.LetItGoConfig) {

	var tplVars map[string]interface{}
	inrec, _ := json.Marshal(cfg)
	json.Unmarshal(inrec, &tplVars)
	inrec, _ = json.Marshal(ctx)
	json.Unmarshal(inrec, &tplVars)

	tplVars["Version"] = cfg.Version()

	if tpl, err := template.New("value").Funcs(sprig.TxtFuncMap()).Parse(*v); err == nil {
		var out bytes.Buffer
		if err := tpl.Execute(&out, tplVars); err == nil {
			*v = out.String()
		}
	}
}
