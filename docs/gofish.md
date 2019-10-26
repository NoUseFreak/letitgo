---
title: "Gofish"
---

Create a PR to update gofish lua scripts.

#### Prerequisites

- Requires `GITHUB_TOKEN` to be set in the environment.

#### Configuration

Parameter | Description | Default
--- | --- | ---
`githubusername` | Github username to for the gofish-food. | ""
`homepage` | Project homepage. | ""
`artifacts` | List of artifact specs. | []
`artifacts[].os` | Os the binary is build for. | ""
`artifacts[].arch` | Architecture the binary is build for. | ""
`artifacts[].url` | Url where the binary can be downloaded. | ""


#### Example

```yaml
   - type: gofish
      githubusername: NoUseFreak
      homepage: https://github.com/NoUseFreak/letitgo
      artifacts:
        - os: darwin
          arch: amd64
          url: https://github.com/NoUseFreak/letitgo/releases/download/{{ .Version }}/darwin_amd64.zip
        - os: linux
          arch: amd64
          url: https://github.com/NoUseFreak/letitgo/releases/download/{{ .Version }}/linux_amd64.zip
```
