---
title: "Continuous integration"
weight: 100
---

To make the build process is documented and versions as much as possible, it is
a good idea to include both build and release scripts in the same repository.

```bash
$ tree ./scripts/
./scripts/
├── build.sh
├── release.sh
└── test.sh
```

Example `scripts/release.sh`.

```
#!/usr/bin/env bash
#
# This script handles the projects release.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Change to the project root.
cd "$DIR"

# Download and run LetItGo
curl -sL http://bit.ly/gh-get | BIN_DIR=/tmp/bin PROJECT=NoUseFreak/letitgo bash
/tmp/bin/letitgo release
```

Make sure the script is executable.
