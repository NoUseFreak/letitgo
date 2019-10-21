---
title: "Archive"
---

#### Configuration

Parameter | Description | Default
--- | --- | ---
`source` | Path to folders that need to be archived. | ""
`target` | Target directory to place archives. | ""
`extras` | List of files to to include in the archive as extras. | []
`method` | What method to use to create the archives. | "zip"

#### Example

```yaml
    - type: archive
      source: "./build/bin/*"
      target: "./build/pkg/"
      extras:
        - "LICENSE"
        - "CHANGELOG.md"
```
