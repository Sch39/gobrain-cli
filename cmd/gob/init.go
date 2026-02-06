package main

import (
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/initcmd"
	"github.com/sch39/gobrain-cli/internal/version"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	opts := &initcmd.Options{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Go project",
		Long:  "Initialize a new Go project in the current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			debug.Printf("init: %v\n", args)
			// Logic to initialize the project
			stableGoVersion, err := version.FetchStableRelease()
			localGo := version.DetectLocalGo()
			if err != nil {
				return err
			}
			debug.Printf("Local Go version: %s\n", localGo)
			debug.Printf("Stable Go Version: %s\n", stableGoVersion)
			debug.Printf("Latest Stable Go Version: %s\n", stableGoVersion[0])
			debug.Printf("Compare semver: %d\n", version.CompareSemver(localGo, stableGoVersion[0]))
			cmd.Println("Project initialized!")
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.Name, "name", "n", "", "Project name")
	f.StringVarP(&opts.Module, "module", "m", "", "Go module path")
	f.StringVarP(&opts.SourceType, "source-type", "t", "", "Source type")
	f.StringVarP(&opts.Source, "source", "s", "", "Source URL")
	f.StringVarP(&opts.Path, "path", "p", "", "Path to project directory")
	f.StringVarP(&opts.Toolchain, "toolchain", "c", "", "Toolchain version")
	f.BoolVar(&opts.Force, "force", false, "Overwrite existing files")

	return cmd
}
