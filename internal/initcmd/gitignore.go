package initcmd

import (
	"os"
	"path/filepath"
	"strings"
)

func EnsureGitignoreManaged(root string) error {
	p := filepath.Join(root, ".gitignore")
	var existing []byte
	if b, err := os.ReadFile(p); err == nil {
		existing = b
	}
	merged := mergeGitignoreContents(existing, nil)
	return os.WriteFile(p, merged, 0o644)
}

func MergeGitignoreFiles(root, templatePath string) error {
	var base []byte
	if b, err := os.ReadFile(filepath.Join(root, ".gitignore")); err == nil {
		base = b
	}
	var extra []byte
	if templatePath != "" {
		if b, err := os.ReadFile(templatePath); err == nil {
			extra = b
		}
	}
	if len(base) == 0 && len(extra) == 0 {
		return nil
	}
	merged := mergeGitignoreContents(base, extra)
	return os.WriteFile(filepath.Join(root, ".gitignore"), merged, 0o644)
}

func mergeGitignoreContents(base, extra []byte) []byte {
	lines := mergeUniqueLines(splitLines(string(base)), splitLines(string(extra)))
	lines = ensureGobIgnored(lines)
	out := strings.Join(lines, "\n")
	if out == "" {
		out = ".gob/"
	}
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	return []byte(out)
}

func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

func mergeUniqueLines(base, extra []string) []string {
	out := make([]string, 0, len(base)+len(extra))
	seen := map[string]bool{}
	for _, ln := range base {
		key := strings.TrimSpace(ln)
		if key != "" && seen[key] {
			continue
		}
		out = append(out, ln)
		if key != "" {
			seen[key] = true
		}
	}
	for _, ln := range extra {
		key := strings.TrimSpace(ln)
		if key == "" || seen[key] {
			continue
		}
		out = append(out, ln)
		seen[key] = true
	}
	return out
}

func ensureGobIgnored(lines []string) []string {
	begin := -1
	end := -1
	for i, ln := range lines {
		switch strings.TrimSpace(ln) {
		case "# === BEGIN GOBRAIN MANAGED ===":
			begin = i
		case "# === END GOBRAIN MANAGED ===":
			end = i
		case ".gob/", ".gob":
			if begin == -1 || end == -1 || i < begin || i > end {
				return lines
			}
		}
	}
	if begin != -1 && end != -1 && begin < end {
		for i := begin + 1; i < end; i++ {
			switch strings.TrimSpace(lines[i]) {
			case ".gob/", ".gob":
				return lines
			}
		}
		out := make([]string, 0, len(lines)+1)
		out = append(out, lines[:begin+1]...)
		out = append(out, ".gob/")
		out = append(out, lines[begin+1:]...)
		return out
	}
	block := []string{
		"# === BEGIN GOBRAIN MANAGED ===",
		".gob/",
		"# === END GOBRAIN MANAGED ===",
	}
	return append(lines, block...)
}
