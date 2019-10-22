---
title: "Helm"
---

Package and/or publish helm charts to a registry like chartmuseum.

#### Prerequisites

- Requires `helm` to be installed.

#### Configuration

Parameter | Description | Default
--- | --- | ---
`charts` | Charts to package. | []
`builddir` | Directory to package charts in. | "./build/helm"
`repository` | Repository to publish charts to. | ""

#### Example

```yaml
    - type: helm:
      charts:
        - ./test/test-chart
      builddir: build/helm/
      repository: https://enxmp9berw1vj.x.pipedream.net
```
