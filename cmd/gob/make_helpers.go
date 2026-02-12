package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/exec"
	"github.com/sch39/gobrain-cli/internal/project"
)

func resolveProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if r, err := project.Find(dir); err == nil {
		debug.Printf("project root found: %s\n", r.Path)
		return r.Path, nil
	}
	debug.Printf("project root default: %s\n", dir)
	return dir, nil
}

func loadConfig(root string) (*config.Config, error) {
	return config.Load(root)
}

func findGenerator(cfg *config.Config, key string) (config.GeneratorCommand, error) {
	name := strings.TrimSpace(key)
	gen, ok := cfg.Generators.Commands[name]
	if !ok {
		return config.GeneratorCommand{}, errors.New("generator command not found: " + name)
	}
	return gen, nil
}

func listGeneratorKeys(cfg *config.Config) []string {
	keys := make([]string, 0, len(cfg.Generators.Commands))
	for k := range cfg.Generators.Commands {
		keys = append(keys, k)
	}
	return keys
}

func formatGeneratorArgs(args []string) string {
	clean := make([]string, 0, len(args))
	for _, a := range args {
		s := strings.TrimSpace(a)
		if s != "" {
			clean = append(clean, s)
		}
	}
	if len(clean) == 0 {
		return "[]"
	}
	return "[" + strings.Join(clean, " ") + "]"
}

func buildTemplateData(cfg *config.Config, argNames []string, argValues map[string]string) map[string]string {
	data := make(map[string]string)
	for _, name := range argNames {
		n := strings.TrimSpace(name)
		if n == "" {
			continue
		}
		u := strings.ToUpper(n[:1]) + n[1:]
		data[u] = argValues[n]
	}
	data["ProjectName"] = strings.TrimSpace(cfg.Project.Name)
	data["ModuleName"] = strings.TrimSpace(cfg.Project.Module)
	return data
}

func ensureModDir(root string) error {
	return os.MkdirAll(filepath.Join(root, ".gob", "mod"), 0o755)
}

func installGeneratorRequires(ctx context.Context, root string, requires []string) error {
	if len(requires) == 0 {
		debug.Printf("make requires: none\n")
		return nil
	}
	r := &exec.OSRunner{}
	for _, dep := range requires {
		d := strings.TrimSpace(dep)
		if d == "" {
			continue
		}
		debug.Printf("make requires: %s\n", d)
		envs := map[string]string{
			"GOMODCACHE":  filepath.Join(root, ".gob", "mod"),
			"GOTOOLCHAIN": "auto",
			"PATH":        filepath.Join(root, ".gob", "bin") + string(os.PathListSeparator) + os.Getenv("PATH"),
		}
		if err := r.Run(ctx, "go", []string{"get", d}, exec.Options{Dir: root, Env: envs}); err != nil {
			return err
		}
	}
	return nil
}
