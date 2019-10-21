---
title: Docker
---

Run the latest version and mount the current directory.

```bash
$ docker pull nousefreak/letitgo
$ docker run -v $(pwd):/app nousefreak/letitgo --version
```
