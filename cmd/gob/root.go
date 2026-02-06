package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gob",
	Short: "Gobrain CLI Tool - Project-scoped Go development CLI",
	Long:  "GoBrain is an project-scoped Go development CLI. It provides a set of tools to help you develop Go projects in a project-scoped environment.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		root, _ := os.Getwd()
		cmd.Println("Current working directory:", root)
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
