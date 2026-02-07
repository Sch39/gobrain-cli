package presets

import _ "embed"

//go:embed presets.yaml
var embeddedPresets string

func DefaultPresetsYAML() string {
	return embeddedPresets
}
