package main

import "github.com/spf13/cobra"

func init() {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Go project",
		Long:  "Initialize a new Go project in the current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Logika init kamu di sini
			cmd.Println("Project initialized!")
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
