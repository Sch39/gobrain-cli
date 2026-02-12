package main

import (
	"context"
	"os"
	"strings"

	"github.com/sch39/gobrain-cli/internal/config"
	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/exec"
	"github.com/spf13/cobra"
)

func NewVerifyCommand() *cobra.Command {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Project verification pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			cwd, _ := os.Getwd()
			debug.Printf("verify: cwd=%s\n", cwd)
			cfg, err := config.Load(cwd)
			if err != nil {
				return err
			}
			envs := baseEnv(cwd)
			r := &exec.OSRunner{}
			for _, step := range cfg.Verify.Pipeline {
				run := strings.TrimSpace(step.Run)
				if run == "" {
					continue
				}
				debug.Printf("verify: step=%s\n", run)
				for _, seg := range splitChain(run) {
					parts := parseCmdLine(seg)
					if len(parts) == 0 {
						continue
					}
					c := parts[0]
					as := parts[1:]
					debug.Printf("verify: cmd=%s args=%v\n", c, as)
					if err := r.Run(ctx, c, as, exec.Options{Dir: cwd, Env: envs}); err != nil {
						if cfg.Verify.FailFast {
							return err
						}
						break
					}
				}
			}
			return nil
		},
	}
	return verifyCmd
}
