package initcmd

import (
	"os"
	"path/filepath"
	"strings"
)

func DetectGoToolchain(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(b), "\n")
	for _, ln := range lines {
		s := strings.TrimSpace(ln)
		if strings.HasPrefix(s, "toolchain ") {
			v := strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(s, "toolchain")), "go")
			if v != "" {
				return v, nil
			}
		}
		if strings.HasPrefix(s, "go ") {
			v := strings.TrimSpace(strings.TrimPrefix(s, "go "))
			if v != "" {
				return v, nil
			}
		}
	}
	return "", nil
}

func DetectGoDirective(root string) (string, error) {
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(b), "\n")
	for _, ln := range lines {
		s := strings.TrimSpace(ln)
		if strings.HasPrefix(s, "go ") {
			v := strings.TrimSpace(strings.TrimPrefix(s, "go "))
			if v != "" {
				return v, nil
			}
		}
	}
	return "", nil
}

func parseToolchain(v string) (int, int, int) {
	parts := strings.Split(strings.TrimSpace(v), ".")
	var maj, min, patch int
	if len(parts) > 0 {
		maj = atoiSafe(parts[0])
	}
	if len(parts) > 1 {
		min = atoiSafe(parts[1])
	}
	if len(parts) > 2 {
		patch = atoiSafe(parts[2])
	}
	return maj, min, patch
}

func parseGoVersion(v string) (int, int) {
	parts := strings.Split(strings.TrimSpace(v), ".")
	var maj, min int
	if len(parts) > 0 {
		maj = atoiSafe(parts[0])
	}
	if len(parts) > 1 {
		min = atoiSafe(parts[1])
	}
	return maj, min
}

func atoiSafe(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			break
		}
		n = n*10 + int(ch-'0')
	}
	return n
}

func ValidateToolchainVsGo(toolchain, goVer string) error {
	if strings.TrimSpace(toolchain) == "" || strings.TrimSpace(goVer) == "" {
		return nil
	}
	tMaj, tMin, _ := parseToolchain(toolchain)
	gMaj, gMin := parseGoVersion(goVer)
	if tMaj < gMaj || (tMaj == gMaj && tMin < gMin) {
		return NewErrToolchainTooLow(toolchain, goVer)
	}
	return nil
}

type ErrToolchainTooLow string

func (e ErrToolchainTooLow) Error() string {
	parts := strings.Split(string(e), "|")
	toolchain := ""
	goVer := ""
	if len(parts) == 2 {
		toolchain = parts[0]
		goVer = parts[1]
	}
	return "toolchain must be >= go directive (go " + goVer + "): got toolchain " + toolchain + " -- choose higher --toolchain or adjust go directive"
}

func NewErrToolchainTooLow(toolchain, goVer string) error {
	return ErrToolchainTooLow(toolchain + "|" + goVer)
}

func EnsureToolchain(root, version string) error {
	p := filepath.Join(root, "go.mod")
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	found := false
	for i, ln := range lines {
		s := strings.TrimSpace(ln)
		if strings.HasPrefix(s, "toolchain ") {
			found = true
			want := "toolchain go" + version
			if s != want {
				lines[i] = want
			}
			break
		}
	}
	if !found {
		lines = append(lines, "toolchain go"+version)
	}
	out := strings.Join(lines, "\n")
	return os.WriteFile(p, []byte(out), 0o644)
}
