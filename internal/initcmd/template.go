package initcmd

import "github.com/sch39/gobrain-cli/internal/presets"

func LoadPresets() ([]presets.Preset, error) {
	return presets.LoadPresets()
}
