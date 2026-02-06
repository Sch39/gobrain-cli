package main

import (
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Go project",
		Long:  "Initialize a new Go project in the current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			debug.Printf("init: %v\n", args)
			// Logika init kamu di sini
			cmd.Println("Project initialized!")
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
