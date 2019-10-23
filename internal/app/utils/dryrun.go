package utils

// DryRun holds dryRun information.
var DryRun dryRun

func init() {
	DryRun = dryRun{}
}

type dryRun struct {
	enabled bool
}

func (d *dryRun) IsEnabled() bool {
	return d.enabled
}

func (d *dryRun) Enable() {
	d.enabled = true
}

func (d *dryRun) Disable() {
	d.enabled = false
}

func (d *dryRun) String() string {
	return map[bool]string{
		true:  "enabled",
		false: "disabled",
	}[d.enabled]
}
