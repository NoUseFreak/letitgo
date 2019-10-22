---
title: "Introduction"
weight: 10
---

LetItGo simplifies automated releases. A simple definition in `.release.yml` in the root of your project is all you need.

{{% block warning %}}
LetItGo prior to version 1.0.0 does not have a stable api.
Since 0.5.0 the config format changed. The goal is to keep it stable. But this is not a guaranty.
{{% /block %}}

### Example

This example shows a few of the features for LetItGo. 

```yaml
letitgo:
  name: letitgo
  description: LetItGo automates releases.
  actions:
    - type: changelog

    - type: githubrelease
      assets:
      - ./build/pkg/*

    - type: homebrew
      homepage: https://github.com/NoUseFreak/letitgo
      url: https://github.com/NoUseFreak/letitgo/releases/download/{{ .Version }}/darwin_amd64.zip
      tap:
        url: git@github.com:NoUseFreak/homebrew-brew.git
      test: system "#{bin}/{{ .Name }} -h"
```