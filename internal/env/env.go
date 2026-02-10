package env

import "os"

func BaseEnv() map[string]string {
	return map[string]string{
		"GOMODCACHE": "./.gob/mod",
		"GOB_BIN":    "./.gob/bin",
	}
}

func Merge(base map[string]string) []string {
	env := os.Environ()
	for k, v := range base {
		env = append(env, k+"="+v)
	}
	return env
}

func WithGobBin() map[string]string {
	return map[string]string{
		"PATH": "./.gob/bin" + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
}
