package main

import (
	"os"
	"path/filepath"
	"strings"
)

func baseEnv(cwd string) map[string]string {
	return map[string]string{
		"GOBIN":       filepath.Join(cwd, ".gob", "bin"),
		"GOMODCACHE":  filepath.Join(cwd, ".gob", "mod"),
		"PATH":        filepath.Join(cwd, ".gob", "bin") + string(os.PathListSeparator) + os.Getenv("PATH"),
		"GOTOOLCHAIN": "auto",
	}
}

func loadDotEnv(p string) map[string]string {
	out := map[string]string{}
	b, err := os.ReadFile(p)
	if err != nil {
		return out
	}
	for _, ln := range strings.Split(string(b), "\n") {
		s := strings.TrimSpace(ln)
		if s == "" {
			continue
		}
		if strings.HasPrefix(s, "#") {
			continue
		}
		i := strings.IndexByte(s, '=')
		if i <= 0 {
			continue
		}
		k := strings.TrimSpace(s[:i])
		v := strings.TrimSpace(s[i+1:])
		out[k] = v
	}
	return out
}

func parseCmdLine(s string) []string {
	var out []string
	var cur strings.Builder
	quote := byte(0)
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if quote != 0 {
			if ch == quote {
				quote = 0
				continue
			}
			cur.WriteByte(ch)
			continue
		}
		if ch == '"' || ch == '\'' {
			quote = ch
			continue
		}
		if ch == ' ' || ch == '\t' {
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteByte(ch)
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

func splitChain(s string) []string {
	parts := strings.Split(s, "&&")
	var out []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
