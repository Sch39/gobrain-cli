package initcmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/exec"
	"go.yaml.in/yaml/v3"
	"golang.org/x/mod/module"
)

const (
	configFileName    = "gob.yaml"
	configFileAltName = "gob.yml"
	configDirName     = ".gob"
)

func Run(ctx context.Context, opts Options) error {
	if opts.Module == "" {
		return errors.New("module name is required")
	}
	if err := module.CheckPath(opts.Module); err != nil {
		return fmt.Errorf("invalid module path: %v", err)
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}
	debug.Printf("init run: root=%s\n", root)
	if configExists(root) && !opts.Force {
		name := existingConfigName(root)
		if name == "" {
			name = configFileName
		}
		return errors.New(name + " already exists")
	}

	hasGoMod := exists(filepath.Join(root, "go.mod"))
	toolchain := ""
	if hasGoMod {
		v, err := DetectGoToolchain(root)
		if err != nil {
			return err
		}
		toolchain = v
		debug.Printf("init run: detected toolchain=%s\n", toolchain)
	} else {
		if strings.TrimSpace(opts.Toolchain) == "" {
			return errors.New("toolchain is required when go.mod is absent")
		}
		toolchain = opts.Toolchain
	}
	debug.Printf("init run: sourceType=%s source=%s\n", opts.SourceType, opts.Source)

	switch strings.ToLower(opts.SourceType) {
	case "preset", "url":
		if strings.TrimSpace(opts.Source) == "" {
			return errors.New("template source is required")
		}
		debug.Printf("init run: prepare from git src=%s path=%s\n", opts.Source, opts.Path)
		if err := PrepareFromGit(ctx, opts.Source, opts.Path, root); err != nil {
			return err
		}
		if !configExists(root) {
			return errors.New("template must contain gob.yaml or gob.yml")
		}
		if exists(filepath.Join(root, "go.mod")) {
			return errors.New("Template cannot contain go.mod")
		}
		if err := HydratePlaceholders(root, opts.Name, opts.Module); err != nil {
			return err
		}
		if cfg, err := config.Load(root); err == nil {
			debug.Printf("init run: layout dirs=%d\n", len(cfg.Layout.Dirs))
			for _, d := range cfg.Layout.Dirs {
				s := strings.TrimSpace(d)
				if s == "" {
					continue
				}
				if err := os.MkdirAll(filepath.Join(root, s), 0o755); err != nil {
					return err
				}
			}
			if !hasGoMod && strings.TrimSpace(toolchain) == "" {
				t := strings.TrimSpace(cfg.Project.Toolchain)
				if t != "" {
					toolchain = t
				}
			}
		}
	case "none", "":
	default:
		return errors.New("invalid source type")
	}

	if !hasGoMod {
		r := &exec.OSRunner{}
		debug.Printf("init run: go mod init %s\n", opts.Module)
		if err := r.Run(ctx, "go", []string{"mod", "init", opts.Module}, exec.Options{Dir: root, Env: map[string]string{"GOTOOLCHAIN": "auto"}}); err != nil {
			return err
		}
		debug.Printf("init run: ensure toolchain %s\n", toolchain)
		if err := EnsureToolchain(root, toolchain); err != nil {
			return err
		}
		if goVer, err := DetectGoDirective(root); err == nil {
			if err := ValidateToolchainVsGo(toolchain, goVer); err != nil {
				return err
			}
		}
	} else {
		if strings.TrimSpace(opts.Toolchain) != "" {
			if goVer, err := DetectGoDirective(root); err == nil {
				if err := ValidateToolchainVsGo(opts.Toolchain, goVer); err != nil {
					return err
				}
			}
			if err := EnsureToolchain(root, opts.Toolchain); err != nil {
				return err
			}
		}
	}

	if err := os.MkdirAll(filepath.Join(root, configDirName), 0o755); err != nil {
		return err
	}
	debug.Printf("init run: ensure gitignore\n")

	if err := EnsureGitignoreManaged(root); err != nil {
		return err
	}

	sourceType := strings.ToLower(opts.SourceType)
	if !configExists(root) || sourceType == "none" {
		debug.Printf("init run: write config new\n")
		type tmplData struct {
			Name      string
			Module    string
			Toolchain string
			Source    string
		}
		data := tmplData{
			Name:      opts.Name,
			Module:    opts.Module,
			Toolchain: toolchain,
			Source:    opts.Source,
		}
		const gobTmpl = `version: "1.0.0"
project:
  name: {{ .Name }}
  module: {{ .Module }}
  toolchain: {{ .Toolchain }}
meta:
  template_origin: {{ .Source }}
  author: ""` + "\n"
		t, err := template.New("gob").Parse(gobTmpl)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(root, configFileName), buf.Bytes()); err != nil {
			return err
		}
	} else if sourceType == "preset" || sourceType == "url" {
		debug.Printf("init run: update config\n")
		cfg, err := config.Load(root)
		if err != nil {
			return err
		}
		cfg.Project.Name = opts.Name
		cfg.Project.Module = opts.Module
		cfg.Project.Toolchain = toolchain
		cfg.Meta.TemplateOrigin = opts.Source
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		path := configPath(root)
		if path == "" {
			return errors.New("config file not found")
		}
		debug.Printf("init run: config path=%s\n", path)
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}
	if err := ensureSingleConfig(root); err != nil {
		return err
	}

	fmt.Println("✔ GoBrain project initialized")
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func configExists(root string) bool {
	if exists(filepath.Join(root, configFileName)) {
		return true
	}
	if exists(filepath.Join(root, configFileAltName)) {
		return true
	}
	return false
}

func configPath(root string) string {
	if exists(filepath.Join(root, configFileName)) {
		return filepath.Join(root, configFileName)
	}
	if exists(filepath.Join(root, configFileAltName)) {
		return filepath.Join(root, configFileAltName)
	}
	return ""
}

func existingConfigName(root string) string {
	if exists(filepath.Join(root, configFileName)) {
		return configFileName
	}
	if exists(filepath.Join(root, configFileAltName)) {
		return configFileAltName
	}
	return ""
}

func ensureSingleConfig(root string) error {
	yamlPath := filepath.Join(root, configFileName)
	ymlPath := filepath.Join(root, configFileAltName)
	yamlExists := exists(yamlPath)
	ymlExists := exists(ymlPath)
	if yamlExists && ymlExists {
		return os.Remove(ymlPath)
	}
	return nil
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
