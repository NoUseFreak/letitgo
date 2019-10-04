package changelog

var changelogTpl = `# Changelog
{{ range .Blocks }}
{{- if gt (len .Commits) 0 -}}
{{ template "block" . }}
{{- end -}}
{{ end }}
`

var blockTpl = `
## {{ .Tag | trim }} ({{ if .Date }}{{ .Date | date "January 2, 2006" }}{{ else }}Unreleased{{ end }})
{{ range .Commits }}
- {{ .Message | trim }} ({{ trunc 8 .Hash.String }})
{{- end }}

`
