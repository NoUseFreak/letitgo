---
title: "Changelog"
---

Generate a changelog and commit it to your project.

#### Prerequisites

- Requires `GITHUB_TOKEN` to be set in the environment.

#### Configuration

Parameter | Description | Default
--- | --- | ---
`file` | File to place the changelog data. | "CHANGELOG.md"
`message` | Commit message to add to the commit. | "Update changelog\n[skip ci]"

#### Example

```yaml
    - type changelog:
      file: CHANGELOG.md
```
