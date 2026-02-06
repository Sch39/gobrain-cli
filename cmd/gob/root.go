package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/project"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gob",
	Short: "Gobrain CLI Tool - Project-scoped Go development CLI",
	Long:  "GoBrain is an project-scoped Go development CLI. It provides a set of tools to help you develop Go projects in a project-scoped environment.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		root, _ := os.Getwd()
		if r, err := project.Find(root); err == nil {
			root = r.Path
		}
		cfg, err := config.Load(root)
		fmt.Printf("root: %v\n", root)
		fmt.Printf("config: %v\n", cfg)
		if err == nil {
			toolChain := strings.TrimSpace(cfg.Project.Toolchain)

			if toolChain != "" {
				currToolchain := strings.TrimSpace(os.Getenv("GOTOOLCHAIN"))
				if currToolchain == "" || currToolchain == "local" {
					_ = os.Setenv("GOTOOLCHAIN", "auto")
				}
			}
		} else {
			fmt.Printf("error loading config: %v\n", err)
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
