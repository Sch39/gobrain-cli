package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/exec"
	"github.com/spf13/cobra"
)

func NewScriptsCommand() *cobra.Command {
	scriptsCmd := &cobra.Command{
		Use:   "scripts",
		Short: "Project scripts runner",
	}

	runCmd := &cobra.Command{
		Use:          "run <name>",
		Short:        "Run a script",
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("need script name. format: gob scripts run <name>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cwd, _ := os.Getwd()
			debug.Printf("scripts run: cwd=%s name=%s\n", cwd, args[0])
			cfg, err := config.Load(cwd)
			if err != nil {
				return err
			}
			name := strings.TrimSpace(args[0])
			seq, ok := cfg.Scripts[name]
			if !ok || len(seq) == 0 {
				return errors.New("script tidak ditemukan atau kosong")
			}
			envs := baseEnv(cwd)
			for k, v := range loadDotEnv(filepath.Join(cwd, ".env")) {
				envs[k] = v
			}
			debug.Printf("scripts run: steps=%d\n", len(seq))
			r := &exec.OSRunner{}
			for _, s := range seq {
				debug.Printf("scripts run: step=%s\n", s)
				for _, seg := range splitChain(s) {
					parts := parseCmdLine(seg)
					if len(parts) == 0 {
						continue
					}
					c := parts[0]
					as := parts[1:]
					if strings.EqualFold(c, "cd") {
						if len(as) < 1 {
							return errors.New("cd membutuhkan path")
						}
						target := strings.TrimSpace(as[0])
						if !filepath.IsAbs(target) {
							target = filepath.Join(cwd, target)
						}
						target = filepath.Clean(target)
						if fi, err := os.Stat(target); err != nil || !fi.IsDir() {
							return errors.New("direktori tidak ditemukan: " + target)
						}
						debug.Printf("scripts run: cd=%s\n", target)
						cwd = target
						continue
					}
					if strings.EqualFold(c, "pwd") {
						debug.Printf("scripts run: pwd=%s\n", cwd)
						os.Stdout.WriteString(cwd + "\n")
						continue
					}
					debug.Printf("scripts run: cmd=%s args=%v\n", c, as)
					if err := r.Run(ctx, c, as, exec.Options{Dir: cwd, Env: envs}); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			debug.Printf("scripts list: cwd=%s\n", cwd)
			cfg, err := config.Load(cwd)
			if err != nil {
				return err
			}
			for k := range cfg.Scripts {
				os.Stdout.WriteString(k + "\n")
			}
			return nil
		},
	}

	scriptsCmd.AddCommand(runCmd)
	scriptsCmd.AddCommand(listCmd)
	return scriptsCmd
}
