package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func parseGeneratorArgs(names []string, raw []string) (map[string]string, error) {
	expected := make([]string, 0, len(names))
	expectedSet := map[string]string{}
	for _, name := range names {
		n := strings.TrimSpace(name)
		if n == "" {
			continue
		}
		expected = append(expected, n)
		expectedSet[strings.ToLower(n)] = n
	}

	values := map[string]string{}
	var positional []string
	for i := 0; i < len(raw); i++ {
		tok := strings.TrimSpace(raw[i])
		if tok == "" {
			continue
		}
		if strings.HasPrefix(tok, "--") {
			key := strings.TrimPrefix(tok, "--")
			if key == "" {
				continue
			}
			val := ""
			if eq := strings.IndexByte(key, '='); eq >= 0 {
				val = key[eq+1:]
				key = key[:eq]
			} else if i+1 < len(raw) && !strings.HasPrefix(raw[i+1], "--") {
				val = raw[i+1]
				i++
			}
			canon := strings.ToLower(strings.TrimSpace(key))
			exp, ok := expectedSet[canon]
			if !ok {
				return nil, fmt.Errorf("unknown generator arg: %s", key)
			}
			values[exp] = strings.TrimSpace(val)
			continue
		}
		positional = append(positional, tok)
	}

	for _, v := range positional {
		for _, name := range expected {
			if _, ok := values[name]; !ok {
				values[name] = v
				break
			}
		}
	}

	var missing []string
	for _, name := range expected {
		if strings.TrimSpace(values[name]) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing generator args. need: %s", strings.Join(missing, " "))
	}
	return values, nil
}

func ensureSnakeCaseFileName(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if name == "" {
		return path
	}
	snake := toSnakeCase(name)
	if snake == "" {
		return path
	}
	return filepath.Join(dir, snake+ext)
}

func toSnakeCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

func toPascalCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 && isAllLower(words[0]) && words[0] == s {
		return s
	}
	var b strings.Builder
	for _, w := range words {
		if w == "" {
			continue
		}
		l := strings.ToLower(w)
		b.WriteString(strings.ToUpper(l[:1]))
		b.WriteString(l[1:])
	}
	return b.String()
}

func toCamelCase(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	if len(words) == 1 && isAllLower(words[0]) && words[0] == s {
		return s
	}
	var b strings.Builder
	for i, w := range words {
		if w == "" {
			continue
		}
		l := strings.ToLower(w)
		if i == 0 {
			b.WriteString(l)
			continue
		}
		b.WriteString(strings.ToUpper(l[:1]))
		b.WriteString(l[1:])
	}
	return b.String()
}

func splitWords(s string) []string {
	var words []string
	var buf strings.Builder
	var prevLower bool
	var prevIsLetter bool
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '_' || ch == '-' || ch == ' ' || ch == '\t' || ch == '.' || ch == '/' || ch == '\\' {
			if buf.Len() > 0 {
				words = append(words, buf.String())
				buf.Reset()
			}
			prevLower = false
			prevIsLetter = false
			continue
		}
		isUpper := ch >= 'A' && ch <= 'Z'
		isLower := ch >= 'a' && ch <= 'z'
		isDigit := ch >= '0' && ch <= '9'
		if buf.Len() > 0 && isUpper && (prevLower || !prevIsLetter) {
			words = append(words, buf.String())
			buf.Reset()
		}
		if isUpper {
			buf.WriteByte(byte(ch))
			prevLower = false
			prevIsLetter = true
			continue
		}
		if isLower || isDigit {
			buf.WriteByte(byte(ch))
			prevLower = isLower
			prevIsLetter = isLower
			continue
		}
		if buf.Len() > 0 {
			words = append(words, buf.String())
			buf.Reset()
		}
		prevLower = false
		prevIsLetter = false
	}
	if buf.Len() > 0 {
		words = append(words, buf.String())
	}
	return words
}

func isAllLower(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < 'a' || ch > 'z' {
			return false
		}
	}
	return true
}
