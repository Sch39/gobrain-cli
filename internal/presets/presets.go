package presets

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

type Preset struct {
	Name   string       `yaml:"name"`
	Source PresetSource `yaml:"source"`
}

type PresetSource struct {
	Repo string `yaml:"repo"`
	Path string `yaml:"path,omitempty"`
}

type presetFile struct {
	Presets []Preset `yaml:"presets"`
}

func LoadPresets() ([]Preset, error) {
	if p := strings.TrimSpace(os.Getenv("GOB_PRESET_FILE")); p != "" {
		if list, err := tryLoadFile(p); err == nil {
			return list, nil
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		if list, ok := loadFromCandidates(cwdCandidates(cwd)); ok {
			return list, nil
		}
	}

	if s := strings.TrimSpace(DefaultPresetsYAML()); s != "" {
		var pf presetFile
		if err := yaml.Unmarshal([]byte(s), &pf); err == nil {
			return sortPresets(pf.Presets), nil
		}
	}
	return []Preset{}, nil
}

func loadFromCandidates(paths []string) ([]Preset, bool) {
	for _, p := range paths {
		if list, err := tryLoadFile(p); err == nil {
			return list, true
		}
	}
	return nil, false
}

func tryLoadFile(path string) ([]Preset, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	list, err := readPresetFile(path)
	if err != nil {
		return nil, err
	}
	return sortPresets(list), nil
}

func cwdCandidates(cwd string) []string {
	return []string{
		filepath.Join(cwd, "presets.yaml"),
		filepath.Join(cwd, "presets.yml"),
		filepath.Join(cwd, ".presets", "presets.yml"),
		filepath.Join(cwd, ".presets", "presets.yaml"),
	}
}

func sortPresets(list []Preset) []Preset {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})
	return list
}

func readPresetFile(path string) ([]Preset, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pf presetFile
	if strings.HasSuffix(strings.ToLower(path), "yaml") || strings.HasSuffix(strings.ToLower(path), "yml") {
		err = yaml.Unmarshal(b, &pf)
	}
	if err != nil {
		return nil, err
	}
	return pf.Presets, nil
}
