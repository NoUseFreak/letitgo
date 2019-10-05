# Github Release

## Prerequisites

- Requires `GITHUB_TOKEN` to be set in the environment.

## Configuration

Parameter | Description | Default
--- | --- | ---
`owner` | Github project owner. | ""
`repo` | Github project repository. | ""
`title` | Title of the release. | "{{ .Version }}"
`description` | Description of the release. | "{{ .Description }}"
`assets` | List of assets to attach to the release. | []

### Example

```yaml
githubrelease:
  - owner: NoUseFreak
    repo: letitgo
    assets:
      - ./build/pkg/*
```
