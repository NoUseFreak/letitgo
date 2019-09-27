package homebrew

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
)

var homebrewTpl = `
class {{ .Name | camelcase }} < Formula
  desc "{{ .Description }}"
  homepage "{{ .Homepage }}"

  version "{{ .Version }}"
  url "{{ .URL }}"
  sha256 "{{ .Hash }}"

{{- with .Dependencies }}
  {{ range $index, $element := . }}
  depends_on "{{ . }}"
  {{- end }}
{{- end -}}

{{- with .Conflicts }}
  {{ range $index, $element := . }}
  conflicts_with "{{ . }}"
  {{- end }}
{{- end }}

  def install
{{ .Install | indent 4 }}
  end

{{ if .Test }}
  test do
{{ .Test | indent 4 }}
  end
{{- end }}
end
`

func (h *Homebrew) templateInput() {
	h.interpolate(&h.URL)
	h.interpolate(&h.Install)
	h.interpolate(&h.Test)
}

func (h *Homebrew) interpolate(v *string) {
	if tpl, err := template.New("value").Funcs(sprig.TxtFuncMap()).Parse(*v); err == nil {
		var out bytes.Buffer
		if err := tpl.Execute(&out, h); err == nil {
			*v = out.String()
		}
	}
}

func (h *Homebrew) template() (string, error) {
	tpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(homebrewTpl)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, h); err != nil {
		return "", err
	}

	return out.String(), nil
}
