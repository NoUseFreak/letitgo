---
title: Example
---

This example `.release.yml` shows a few of the features for LetItGo. 


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

You can now run LetItGo and make it run the actions.

```bash
$ letitgo release
```
