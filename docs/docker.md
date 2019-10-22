---
title: "Docker"
---

Builds, tags and pushes docker images.

#### Prerequisites

- Requires `DOCKER_AUTH_PASSWORD` to be set in the environment.

#### Configuration

Parameter | Description | Default
--- | --- | ---
`images` | List of images. | []
`dockerfile` | Path to Dockerfile. | "Dockerfile"
`nopush` | State if we want to push to a registry | false
`auth.username` | Username of the registry. | ""
`auth.password` | Password of the registry. | ""

#### Example

The following example configuration will update `Formula/letitgo.rb`.

```yaml
    - type: docker
      dockerfile: "./Dockerfile"
      images:
        - "nousefreak/test"
        - "nousefreak/test:{{ .Version }}"
        - "nousefreak/test:{{ .Version.Major }}"
        - "nousefreak/test:{{ .Version.Major }}.{{ .Version.Minor }}"
      auth:
        username: nousefreak
      nopush: false
```
