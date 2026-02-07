package version

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func DetectLocalGo() string {
	out, err := exec.Command("go", "version").Output()
	if err == nil {
		parts := strings.Fields(string(out))
		for _, p := range parts {
			if strings.HasPrefix(p, "go1.") {
				return strings.TrimPrefix(p, "go")
			}
		}
	}
	v := runtime.Version()
	if strings.HasPrefix(v, "go1.") {
		return strings.TrimPrefix(v, "go")
	}
	return "1.23.0"
}

func CompareSemver(a, b string) int {
	ap, bp := strings.Split(a, "."), strings.Split(b, ".")
	fill := func(s []string) []int {
		res := make([]int, 3)
		for i := 0; i < 3; i++ {
			if i < len(s) {
				res[i] = atoi(s[i])
			}
		}
		return res
	}
	va, vb := fill(ap), fill(bp)
	for i := 0; i < 3; i++ {
		if va[i] < vb[i] {
			return -1
		}
		if va[i] > vb[i] {
			return 1
		}
	}
	return 0
}

func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			break
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}

func FetchStableRelease() []string {
	type rel struct {
		Version string `json:"version"`
		Stable  bool   `json:"stable"`
	}
	client := &http.Client{Timeout: 3 * time.Second}
	for _, u := range []string{"https://go.dev/dl/?mode=json", "https://golang.org/dl/?mode=json"} {
		resp, err := client.Get(u)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		var list []rel
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			continue
		}

		var vlist []string
		re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
		for _, r := range list {
			v := strings.TrimPrefix(r.Version, "go")
			if r.Stable && re.MatchString(v) {
				vlist = append(vlist, v)
			}
		}
		if len(vlist) > 0 {
			// Sort descending
			for i := 0; i < len(vlist)-1; i++ {
				for j := i + 1; j < len(vlist); j++ {
					if CompareSemver(vlist[i], vlist[j]) < 0 {
						vlist[i], vlist[j] = vlist[j], vlist[i]
					}
				}
			}
			return vlist
		}
	}
	return []string{}
}
