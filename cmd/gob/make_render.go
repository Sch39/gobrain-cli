package gob

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
)

func renderGeneratorTemplates(root string, templates []config.Template, data map[string]string) error {
	funcs := template.FuncMap{
		"lower":       func(s string) string { return strings.ToLower(s) },
		"upper":       func(s string) string { return strings.ToUpper(s) },
		"title":       func(s string) string { return strings.Title(s) },
		"snake_case":  toSnakeCase,
		"pascal_case": toPascalCase,
		"camel_case":  toCamelCase,
		"pascalcase":  toPascalCase,
		"camelcase":   toCamelCase,
	}
	for _, t := range templates {
		src := filepath.Join(root, filepath.FromSlash(strings.TrimSpace(t.Src)))
		destTmpl := strings.TrimSpace(t.Dest)
		if destTmpl == "" {
			return errors.New("template dest is empty")
		}
		debug.Printf("make template: src=%s dest=%s\n", src, destTmpl)
		dt, err := template.New("dest").Funcs(funcs).Parse(destTmpl)
		if err != nil {
			return err
		}
		var dbuf bytes.Buffer
		if err := dt.Execute(&dbuf, data); err != nil {
			return err
		}
		dest := filepath.Join(root, filepath.FromSlash(dbuf.String()))
		dest = ensureSnakeCaseFileName(dest)
		debug.Printf("make template: out=%s\n", dest)
		b, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		tpl, err := template.New("file").Funcs(funcs).Parse(string(b))
		if err != nil {
			return err
		}
		var out bytes.Buffer
		if err := tpl.Execute(&out, data); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, out.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}
