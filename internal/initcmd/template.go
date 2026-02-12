package initcmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sch39/gobrain-cli/internal/exec"
	"github.com/sch39/gobrain-cli/internal/presets"
)

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func HydratePlaceholders(root, name, module string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			n := info.Name()
			if n == ".git" || n == ".gob" || n == ".template" || n == "vendor" || n == "bin" || n == "build" || n == "dist" || n == ".idea" || n == ".vscode" {
				return filepath.SkipDir
			}
		}
		if info.IsDir() {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(b)
		reProj := regexp.MustCompile(`{{\s*\.ProjectName\s*}}`)
		reMod := regexp.MustCompile(`{{\s*\.ModuleName\s*}}`)
		out := reProj.ReplaceAllString(s, name)
		out = reMod.ReplaceAllString(out, module)
		if out != s {
			if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
				return err
			}
		}
		return nil
	})
}

func PrepareFromGit(ctx context.Context, url, subpath, dst string) error {
	trim := strings.TrimSpace(subpath)
	isRoot := trim == "" || trim == "." || trim == "/" || trim == `\`
	if fi, err := os.Stat(url); err == nil && fi.IsDir() {
		src := url
		if !isRoot {
			sp := filepath.Clean(filepath.FromSlash(trim))
			p := filepath.Join(url, sp)
			rel, err := filepath.Rel(url, p)
			if err != nil {
				return err
			}
			if strings.HasPrefix(rel, "..") {
				return os.ErrPermission
			}
			src = p
			if _, err := os.Stat(src); err != nil {
				return err
			}
		}
		return copyDir(src, dst)
	}
	tmp, err := os.MkdirTemp("", "gob-tpl-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	r := &exec.OSRunner{}
	args := []string{"clone", "--depth", "1", url, tmp}
	token := strings.TrimSpace(os.Getenv("GOB_GIT_TOKEN"))
	opts := exec.Options{}
	if token != "" {
		args = []string{
			"-c", "http.extraHeader=Authorization: Bearer " + token,
			"clone", "--depth", "1", url, tmp,
		}
	}
	if gl := strings.TrimSpace(os.Getenv("GOB_GITLAB_TOKEN")); gl != "" {
		args = append([]string{"-c", "http.extraHeader=Authorization: Bearer " + gl}, args...)
	}
	user := strings.TrimSpace(os.Getenv("GOB_GIT_BASIC_USER"))
	pass := strings.TrimSpace(os.Getenv("GOB_GIT_BASIC_TOKEN"))
	if user != "" && pass != "" {
		cred := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
		args = append([]string{"-c", "http.extraHeader=Authorization: Basic " + cred}, args...)
	}
	if hdrs := strings.TrimSpace(os.Getenv("GOB_GIT_EXTRA_HEADERS")); hdrs != "" {
		sep := []string{"\n", "\r\n", ";", ","}
		for _, s := range sep {
			hdrs = strings.ReplaceAll(hdrs, s, "\n")
		}
		for _, ln := range strings.Split(hdrs, "\n") {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			args = append([]string{"-c", "http.extraHeader=" + ln}, args...)
		}
	}
	if key := strings.TrimSpace(os.Getenv("GOB_GIT_SSH_KEY")); key != "" {
		opts.Env = map[string]string{
			"GIT_SSH_COMMAND": fmt.Sprintf("ssh -i \"%s\" -o IdentitiesOnly=yes", key),
		}
	}
	if err := r.Run(ctx, "git", args, opts); err != nil {
		return err
	}
	src := tmp
	if !isRoot {
		sp := filepath.Clean(filepath.FromSlash(trim))
		p := filepath.Join(tmp, sp)
		rel, err := filepath.Rel(tmp, p)
		if err != nil {
			return err
		}
		if strings.HasPrefix(rel, "..") {
			return os.ErrPermission
		}
		src = p
		if _, err := os.Stat(src); err != nil {
			return err
		}
	}
	return copyDir(src, dst)
}

func LoadPresets() ([]presets.Preset, error) {
	return presets.LoadPresets()
}
