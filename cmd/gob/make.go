package main

import (
	"context"
	"errors"
	"os"
	"sort"

	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/spf13/cobra"
)

func NewMakeCommand() *cobra.Command {
	makeCmd := &cobra.Command{
		Use:   "make",
		Short: "Code generators based on templates",
		Long:  "Run generator by command name. Args can be positional or flags like --name=value.",
		Example: `gob make handler User
gob make handler --name=User`,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, a := range args {
				if a == "-h" || a == "--help" {
					return cmd.Help()
				}
			}
			if len(args) < 1 {
				return errors.New("need generator command. format: gob make <command> [args...]")
			}
			ctx := context.Background()
			root, err := resolveProjectRoot()
			if err != nil {
				return err
			}
			debug.Printf("make run: root=%s\n", root)
			cfg, err := loadConfig(root)
			if err != nil {
				return err
			}
			gen, err := findGenerator(cfg, args[0])
			if err != nil {
				return err
			}
			debug.Printf("make run: cmd=%s args=%v\n", args[0], args[1:])
			argValues, err := parseGeneratorArgs(gen.Args, args[1:])
			if err != nil {
				return err
			}
			data := buildTemplateData(cfg, gen.Args, argValues)
			if err := ensureModDir(root); err != nil {
				return err
			}
			debug.Printf("make run: requires=%d templates=%d\n", len(cfg.Generators.Requires), len(gen.Templates))
			if err := installGeneratorRequires(ctx, root, cfg.Generators.Requires); err != nil {
				return err
			}
			if err := renderGeneratorTemplates(root, gen.Templates, data); err != nil {
				return err
			}
			return nil
		},
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List generators",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveProjectRoot()
			if err != nil {
				return err
			}
			debug.Printf("make list: root=%s\n", root)
			cfg, err := loadConfig(root)
			if err != nil {
				return err
			}
			keys := listGeneratorKeys(cfg)
			sort.Strings(keys)
			for _, k := range keys {
				argsLabel := formatGeneratorArgs(cfg.Generators.Commands[k].Args)
				os.Stdout.WriteString("- " + k + " " + argsLabel + "\n")
			}
			return nil
		},
	}
	makeCmd.AddCommand(listCmd)
	return makeCmd
}
