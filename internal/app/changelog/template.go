package changelog

var changelogTpl = `# Changelog
{{ range .Blocks }}
{{- if gt (len .Commits) 0 -}}
{{ template "block" . }}
{{- end -}}
{{ end }}
`

var blockTpl = `
## {{ .Tag | trim }} ({{ if .Date }}{{ dateInZone "January 2, 2006" .Date "UTC" }}{{ else }}Unreleased{{ end }})
{{ range .Commits }}
- {{ .Message | trim }} ({{ trunc 8 .Hash.String }})
{{- end }}

`
