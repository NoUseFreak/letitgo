package homebrew

var homebrewTpl = `# This file was generated by LetItGo.
class {{ camelcase (snakecase .Name) }} < Formula
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
