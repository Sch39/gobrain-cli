package gob

import (
	"context"
	"errors"
	"os"

	"github.com/sch39/gobrain-cli/internal/debug"
	"github.com/sch39/gobrain-cli/internal/exec"
	"github.com/spf13/cobra"
)

func NewExecCommand() *cobra.Command {
	execCmd := &cobra.Command{
		Use:               "exec [command] [args...]",
		Short:             "Execute a command with project-scoped env",
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, a := range args {
				if a == "-h" || a == "--help" {
					return cmd.Help()
				}
			}
			if len(args) < 1 {
				return errors.New("need command to execute")
			}
			ctx := context.Background()
			cwd, _ := os.Getwd()
			root, err := resolveProjectRoot()
			if err != nil {
				return err
			}
			debug.Printf("exec: root=%s cwd=%s cmd=%s args=%v\n", root, cwd, args[0], args[1:])
			envs := baseEnv(root)
			r := &exec.OSRunner{}
			return r.Run(ctx, args[0], args[1:], exec.Options{Dir: cwd, Env: envs})
		},
	}
	return execCmd
}
