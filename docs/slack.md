---
title: "Slack"
---

Publish a success message when the release completes.

#### Prerequisites

- Requires `SLACK_TOKEN` to be set in the environment.

#### Configuration

Parameter | Description | Default
--- | --- | ---
`channel` | Channel to post to. | "#released"
`message` | Post message. | "Project release completed"
`author` | Post author. | "LetItGo"

#### Example

The following example configuration will post a message to `#released`.

```yaml
    - type: slack
      channel: "#released"
```
