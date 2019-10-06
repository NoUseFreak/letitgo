# Snapcraft

## Prerequisites

- Requires `snapcraft` to be installed and logged in.

## Configuration

Parameter | Description | Default
--- | --- | ---
`assets` | Binary to package in the snap. | []
`architecture` | Type of binary. (all, amd64, i386) | ""

### Example

```yaml
snapcraft:
  - assets: 
      - build/bin/linux_amd64/letitgo
    architecture: amd64
```
