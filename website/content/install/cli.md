---
title: CLI
---

Installs the latest release to `/usr/local/bin`.

```bash
$ curl -sL http://bit.ly/gh-get | PROJECT=NoUseFreak/letitgo bash
$ letitgo --version
```

If you don't have write permissions to `/usr/local/bin`. You can overwrite the
`BIN_DIR`.

```bash
$ curl -sL http://bit.ly/gh-get | BIN_DIR=/tmp/bin PROJECT=NoUseFreak/letitgo bash
$ /tmp/bin/letitgo --version
```