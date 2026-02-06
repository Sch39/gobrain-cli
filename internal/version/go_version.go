package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

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

func FetchStableRelease() ([]string, error) {
	urls := []string{
		"https://go.dev/dl/?mode=json",
		"https://golang.org/dl/?mode=json",
	}
	client := &http.Client{Timeout: 3 * time.Second}

	for _, u := range urls {
		resp, err := client.Get(u)

		if err != nil {
			continue
		}

		defer resp.Body.Close()
		var releases []GoVersion
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			continue
		}

		var versions []string
		re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
		for _, r := range releases {
			v := strings.TrimPrefix(r.Version, "go")
			if r.Stable && re.MatchString(v) {
				versions = append(versions, v)
			}
		}

		sort.Slice(versions, func(i, j int) bool {
			return CompareSemver(versions[i], versions[j]) > 0
		})
		return versions, nil
	}
	return nil, fmt.Errorf("failed to fetch stable release")
}

func CompareSemver(a string, b string) int {
	parse := func(v string) []int {
		parts := strings.Split(v, ".")
		res := make([]int, 3)
		for i := 0; i < len(parts) && i < 3; i++ {
			fmt.Sscanf(parts[i], "%d", &res[i])
		}
		return res
	}
	va, vb := parse(a), parse(b)
	for i := 0; i < 3; i++ {
		if va[i] != vb[i] {
			return va[i] - vb[i]
		}
	}
	return 0
}
