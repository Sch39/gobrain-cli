package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/initcmd"
	"github.com/sch39/gobrain-cli/internal/version"
	"github.com/spf13/cobra"
	mod "golang.org/x/mod/module"
)

func NewInitCommand() *cobra.Command {
	var opts initcmd.Options

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Go project",
		Long:  "Initialize a new Go project in the current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			debug.Printf("init start")
			if !opts.Force {
				for _, f := range []string{"gob.yaml", "gob.yml"} {
					if _, err := os.Stat(f); err == nil {
						return fmt.Errorf("%s already exists: project already initialized", f)
					}
				}
			}

			// interactive prompts
			if err := runInitPrompts(&opts, cmd); err != nil {
				return err
			}
			debug.Println("Project initialized!")
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

func runInitPrompts(opts *initcmd.Options, cmd *cobra.Command) error {
	if opts.Name == "" {
		if err := survey.AskOne(&survey.Input{Message: "Project name:"}, &opts.Name, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	if opts.Module == "" {
		prompt := &survey.Input{Message: "What is your module path?", Help: "e.g., github.com/user/project"}
		if err := survey.AskOne(prompt, &opts.Module, survey.WithValidator(func(val interface{}) error {
			s, ok := val.(string)
			if !ok {
				return errors.New("invalid module path")
			}
			return mod.CheckPath(s)
		}), survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	if !cmd.Flags().Changed("source-type") {
		survey.AskOne(&survey.Select{
			Message: "Template source:",
			Options: []string{"none", "preset", "url"},
			Default: "none",
		}, &opts.SourceType)
	}

	// Source Logic
	switch opts.SourceType {
	case "preset":
		if opts.Source == "" {
			list, err := initcmd.LoadPresets()
			if err != nil {
				return err
			}
			var names []string
			for _, p := range list {
				names = append(names, p.Name)
			}
			if len(names) == 0 {
				return errors.New("preset list kosong")
			}

			var chosen string
			if err := survey.AskOne(&survey.Select{Message: "Choose preset:", Options: names}, &chosen); err != nil {
				return err
			}
			for _, p := range list {
				if p.Name == chosen {
					opts.Source = p.Source.Repo
					break
				}
			}
		}
	case "url":
		if opts.Source == "" {
			survey.AskOne(&survey.Input{Message: "Template URL/path:"}, &opts.Source, survey.WithValidator(survey.Required))
		}
	}

	// Toolchain Logic
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) && opts.Toolchain == "" {
		local := version.DetectLocalGo()
		debug.Printf("Detected local Go version: Go %s\n", local)

		releases := version.FetchStableRelease()
		optsList := []string{"go" + local + " (local)"}
		for _, v := range releases {
			if version.CompareSemver(v, local) > 0 {
				optsList = append(optsList, "go"+v)
			}
		}
		optsList = append(optsList, "Custom...")

		var selected string
		if err := survey.AskOne(&survey.Select{Message: "Go toolchain:", Options: optsList, Default: optsList[0]}, &selected); err != nil {
			return err
		}

		if selected == "Custom..." {
			survey.AskOne(&survey.Input{Message: "Enter Go toolchain (Major.Minor.Patch):", Default: local}, &opts.Toolchain, survey.WithValidator(survey.Required))
		} else {
			opts.Toolchain = strings.TrimPrefix(strings.Split(selected, " ")[0], "go")
		}

		re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
		if !re.MatchString(opts.Toolchain) {
			return errors.New("invalid toolchain format")
		}
		if version.CompareSemver(opts.Toolchain, local) < 0 {
			return errors.New("toolchain must be >= local version")
		}
	}

	return nil
}
