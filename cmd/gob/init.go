package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/initcmd"
	"github.com/sch39/gobrain-cli/internal/version"
	"github.com/spf13/cobra"
	mod "golang.org/x/mod/module"
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
			stableGoVersion := version.FetchStableRelease()
			localGo := version.DetectLocalGo()
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

func validateEnv(force bool) error {
	if !force {
		for _, f := range []string{"gob.yaml", "gob.yml"} {
			if _, err := os.Stat(f); err == nil {
				return fmt.Errorf("file %s already exists", f)
			}
		}
	}
	return nil
}

func promptMissingInputs(opts *initcmd.Options, cmd *cobra.Command) error {
	if opts.Name == "" {
		survey.AskOne(&survey.Input{Message: "Project name:"}, &opts.Name, survey.WithValidator(survey.Required))
	}

	if opts.Module == "" {
		survey.AskOne(&survey.Input{Message: "Module path:"}, &opts.Module, survey.WithValidator(func(val interface{}) error {
			return mod.CheckPath(val.(string))
		}))
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) && opts.Toolchain == "" {

	}

	return nil
}
