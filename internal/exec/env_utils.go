package exec

import "strings"

func mergeEnv(base []string, override map[string]string) []string {
	envMap := make(map[string]string, len(base))

	for _, e := range base {
		if kv := splitEnv(e); kv != nil {
			envMap[kv[0]] = kv[1]
		}
	}

	for k, v := range override {
		envMap[k] = v
	}

	result := make([]string, 0, len(envMap))
	for k, v := range envMap {
		result = append(result, k+"="+v)
	}

	return result
}

func splitEnv(e string) []string {
	for i := 0; i < len(e); i++ {
		if e[i] == '=' {
			return []string{e[:i], e[i+1:]}
		}
	}
	return nil
}

func sanitizeArgs(args []string) []string {
	safe := make([]string, len(args))

	for i, a := range args {
		lower := strings.ToLower(a)

		switch {
		case strings.Contains(lower, "authorization:"):
			safe[i] = "AUTHORIZATION: *****"
		case strings.Contains(lower, "private-token:"):
			safe[i] = "PRIVATE-TOKEN: *****"
		case strings.HasPrefix(a, "http") && strings.Contains(a, "@"):
			safe[i] = "URL(REDACTED)"
		default:
			safe[i] = a
		}
	}

	return safe
}
