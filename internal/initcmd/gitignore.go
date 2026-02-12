package initcmd

import (
	"bytes"
	"os"
	"path/filepath"
)

func EnsureGitignoreManaged(root string) error {
	p := filepath.Join(root, ".gitignore")
	block := []byte("\n# === BEGIN GOBRAIN MANAGED ===\n.gob/\n# === END GOBRAIN MANAGED ===\n")
	if _, err := os.Stat(p); err != nil {
		return os.WriteFile(p, block, 0o644)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	if bytes.Contains(b, block) {
		return nil
	}
	b = append(b, block...)
	return os.WriteFile(p, b, 0o644)
}
