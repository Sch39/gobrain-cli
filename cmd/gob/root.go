package main

import (
	"os"
	"strings"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/project"
	"github.com/spf13/cobra"
)

var debugFlag bool

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

		if debugFlag || (err == nil && cfg.Debug) {
			debug.Set(true)
		} else {
			debug.Set(false)
		}
		debug.Printf("root: %v\n", root)
		if err == nil {
			toolChain := strings.TrimSpace(cfg.Project.Toolchain)

			if toolChain != "" {
				currToolchain := strings.TrimSpace(os.Getenv("GOTOOLCHAIN"))
				if currToolchain == "" || currToolchain == "local" {
					_ = os.Setenv("GOTOOLCHAIN", "auto")
				}
			}
		} else {
			debug.Printf("error loading config: %v\n", err)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug output")
}

func Execute() {
	// Add commands
	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewMakeCommand())
	rootCmd.AddCommand(NewExecCommand())
	rootCmd.AddCommand(NewScriptsCommand())
	rootCmd.AddCommand(NewToolsCommand())
	rootCmd.AddCommand(NewVerifyCommand())

	// Execute the root command
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
