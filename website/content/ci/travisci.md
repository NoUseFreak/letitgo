---
title: Travis CI
---

#### Setup

Example `.travisci.yml` file.


```yaml
# build steps
language: go
go:
- stable

# deploy steps
deploy:
- provider: script
  skip_cleanup: true # keeps any build artifacts around
  script: bash scripts/release.sh
  on:
    repo: NoUseFreak/letitgo
    branch: master
    tags: true
```

#### Secrets

If you want to add secrets to the build environment, travis-ci provides a 
[CLI tool](https://docs.travis-ci.com/user/encryption-keys/#usage) to do this.

```
$ travis encrypt GITHUB_TOKEN=<token> --add
```
