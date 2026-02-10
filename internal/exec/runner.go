package exec

import "context"

type Runner interface {
	Run(ctx context.Context, cmd string, args []string, opts Options) error
}
