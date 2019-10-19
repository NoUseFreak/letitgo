# Homebrew

Currently it is only supported to update [Taps](https://docs.brew.sh/Taps).
It requires `GITHUB_TOKEN` to be set. 

## Prerequisites

- Requires `GITHUB_TOKEN` to be set in the environment.

## Configuration

Parameter | Description | Default
--- | --- | ---
`name` | Name of the package. | ""
`description` | Description of the package. | ""
`url` | Link to binary that needs to be added. | ""
`version` | Version of the package | "{{ .Version }}"
`dependencies` | List of dependencies. | []
`conflicts` | List of conflicting packages. | []
`tap.url` | Git repository to publish to. | ""
`folder` | Folder to place the homebrew spec. | "Formula"
`install` | Install action | "bin.install \"{{ .Name }}\""
`test` | Test action | ""

### Example

The following example configuration will update `Formula/letitgo.rb`.

```yaml
    - type: homebrew:
      homepage: https://github.com/NoUseFreak/letitgo
      url: https://github.com/NoUseFreak/letitgo/releases/download/{{ .Version }}/darwin_amd64.zip
      tap:
        url: git@github.com:NoUseFreak/homebrew-brew.git
      test: system "#{bin}/{{ .Name }} -h"
```
