package gofish

var luaTpl = `-- This file was generated by LetItGo.
local name = "{{ .Name }}"
local version = "{{ .Version }}"

food = {
    name = name,
    description = "{{ .Description }}",
    homepage = "{{ .Homepage }}",
    version = version,
    packages = {
		{{- range $i, $value := .Artifacts -}}
		{{ if ne $i 0 -}},{{- end }}
        {
            os = "{{ $value.Os }}",
            arch = "{{ $value.Arch }}",
            url = "{{ $value.URL }}",
            sha256 = "{{ $value.Sha256 }}",
            resources = {
                {
                    path = name,
                    installpath = "bin/" .. name,
                    executable = true
                }
            }
        }
		{{- end }}
    }
}

`
