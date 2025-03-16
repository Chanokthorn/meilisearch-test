package runner

import "context"

type Runner interface {
	Run(ctx context.Context) (finalTaskUID int, err error)
}
