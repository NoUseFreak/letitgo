# LetItGo

[![Build status](https://img.shields.io/travis/NoUseFreak/letitgo/master?style=flat-square)](https://travis-ci.org/NoUseFreak/letitgo)
[![Release](https://img.shields.io/github/v/release/NoUseFreak/letitgo?style=flat-square)](https://github.com/NoUseFreak/letitgo/releases)
[![Maintained](https://img.shields.io/maintenance/yes/2019?style=flat-square)](https://github.com/NoUseFreak/letitgo)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/nousefreak/letitgo?style=flat-square)](https://hub.docker.com/r/nousefreak/letitgo)
![Website](https://img.shields.io/netlify/7c9a64af-aefa-4157-b681-20833d61f7d1?style=flat-square)
[![License](https://img.shields.io/github/license/NoUseFreak/letitgo?style=flat-square)](https://github.com/NoUseFreak/letitgo/blob/master/LICENSE)
[![Coffee](https://img.shields.io/badge/☕️-Buy%20me%20a%20coffee-blue?style=flat-square&color=blueviolet)](https://www.buymeacoffee.com/driesdepeuter)

LetItGo simplifies automated releases. A simple definition in `.release.yml` in 
the root of your project is all you need.

Check out the [docs](https://letitgo.nousefreak.be/) for full documentation. 


## Install

__Homebrew__

```bash
brew install NoUseFreak/brew/letitgo
letitgo --version
```

__CLI__

```bash
curl -sL http://bit.ly/gh-get | PROJECT=NoUseFreak/letitgo bash
letitgo --version
```

__Docker__

```bash
docker run -v $(pwd):/app nousefreak/letitgo --version
```

__Anywhere__

```bash
curl -sL http://bit.ly/gh-get | BIN_DIR=/tmp/bin PROJECT=NoUseFreak/letitgo bash
/tmp/bin/letitgo --version
```

## Usage

The most common use case would be to provide a `.release.yml` file in the root
of your project, and let `letitgo` do it's thing.

```bash
$ letitgo
```

## Init

You can use `letitgo init` to help you generate your `.release.yml` file.

It will go through all available actions and provide you with an example for
each of the actions. 

## Actions

Actions as as the name explains, actions that letitgo need to execute when
the release process is triggered.

Action | Description
--- | ---
archive | Create archives of files.
changelog | Generate a changelog and commit it to your project.
githubrelease | Publish generated artifacts and attach them to a github release.
helm | Package and/or publish helm charts to a registry like chartmuseum.
homebrew | Update your personal homebrew tap with your latest config.
snapcraft | Package and upload your snap to snapcraft.

All actions and example configuration can be found in the [docs directory](docs/).

## Example

The following is an example config used to release this project.

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
