package main

import "github.com/spf13/cobra"

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show GoBrain version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			printVersionUI()
			return nil
		},
	}
}
