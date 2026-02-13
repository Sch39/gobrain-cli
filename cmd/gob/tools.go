package gob

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/exec"
	"github.com/spf13/cobra"
)

func NewToolsCommand() *cobra.Command {
	toolsCmd := &cobra.Command{
		Use:   "tools",
		Short: "Project-scoped tooling",
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install declared tools to ./.gob/bin",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cwd, _ := os.Getwd()
			debug.Printf("tools install: cwd=%s\n", cwd)
			if err := os.MkdirAll(filepath.Join(cwd, ".gob", "bin"), 0o755); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Join(cwd, ".gob", "mod"), 0o755); err != nil {
				return err
			}
			cfg, err := config.Load(cwd)
			if err != nil {
				return err
			}
			if len(cfg.Tools) == 0 {
				return errors.New("no tools declared in gob.yaml")
			}
			r := &exec.OSRunner{}
			for _, t := range cfg.Tools {
				pkg := t.Pkg
				if pkg == "" {
					continue
				}
				debug.Printf("tools install: %s\n", pkg)
				envs := baseEnv(cwd)
				if err := r.Run(ctx, "go", []string{"install", pkg}, exec.Options{Dir: cwd, Env: envs}); err != nil {
					return err
				}
			}
			fmt.Println("✔ Tools installed to ./.gob/bin")
			return nil
		},
	}

	runCmd := &cobra.Command{
		Use:   "run [tool] [args...]",
		Short: "Run a tool with project-scoped PATH injection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("need tool name. format: gob tools run <tool> [args...]")
			}
			return nil
		},
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cwd, _ := os.Getwd()
			debug.Printf("tools run: cwd=%s args=%v\n", cwd, args)
			bin := args[0]
			if _, err := os.Stat(filepath.Join(cwd, ".gob", "bin", bin+".exe")); err == nil {
				bin = bin + ".exe"
			}
			envs := baseEnv(cwd)
			r := &exec.OSRunner{}
			debug.Printf("tools run: bin=%s\n", bin)
			return r.Run(ctx, filepath.Join(cwd, ".gob", "bin", bin), args[1:], exec.Options{Dir: cwd, Env: envs})
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List project tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			debug.Printf("tools list: cwd=%s\n", cwd)
			cfg, err := config.Load(cwd)
			if err != nil {
				return err
			}
			for _, t := range cfg.Tools {
				name := t.Name
				bin := name
				if runtime.GOOS == "windows" && filepath.Ext(bin) == "" {
					bin = bin + ".exe"
				}
				p := filepath.Join(cwd, ".gob", "bin", bin)
				if _, err := os.Stat(p); err == nil {
					fmt.Printf("%s installed\n", name)
				} else {
					fmt.Printf("%s missing\n", name)
				}
			}
			return nil
		},
	}

	toolsCmd.AddCommand(installCmd)
	toolsCmd.AddCommand(runCmd)
	toolsCmd.AddCommand(listCmd)
	return toolsCmd
}
