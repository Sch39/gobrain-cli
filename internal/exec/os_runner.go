package exec

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type OsRunner struct{}

type OSRunner struct{}

func (r *OSRunner) Run(
	ctx context.Context,
	command string,
	args []string,
	opts Options,
) error {
	cmd := buildCommand(ctx, command, args, opts)

	if err := cmd.Run(); err != nil {
		safeArgs := sanitizeArgs(args)
		return fmt.Errorf("command failed: %s %v: %w", command, safeArgs, err)
	}

	return nil
}

func buildCommand(
	ctx context.Context,
	command string,
	args []string,
	opts Options,
) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)

	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = resolveEnv(opts.Env)

	return cmd
}

func resolveEnv(override map[string]string) []string {
	if len(override) == 0 {
		return os.Environ()
	}
	return mergeEnv(os.Environ(), override)
}
