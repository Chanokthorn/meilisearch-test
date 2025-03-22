package runner

import (
	"context"
	"ms-tester/storage"
)

type Runner interface {
	Run(ctx context.Context, loader storage.StreamLoader) (finalTaskUID int, err error)
}
