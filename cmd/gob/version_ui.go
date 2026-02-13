package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/sch39/gobrain-cli/internal/ui"
)

var appVersion = "v0.3.0"

const appDescription = "Project-scoped Go development CLI"

func printVersionUI() {
	art := strings.TrimRight(ui.BrandArt, "\n")
	if art != "" {
		fmt.Println(art)
	}
	fmt.Printf("Version: %s\n", appVersion)
	fmt.Println(appDescription)
}
