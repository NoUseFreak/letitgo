package snapcraft

var snapcraftTpl = `
name: {{ default .LetItGo.Name .Name }}
version: {{ .Version }}
summary: {{ default .LetItGo.Description .Description }}
description: {{ default .LetItGo.Description .Description }}

confinement: strict
base: core18
architectures: 
  - {{ .Architecture }}
apps:
  {{ default .LetItGo.Name .Name }}:
    command: {{ default .LetItGo.Name .Name }}
`
