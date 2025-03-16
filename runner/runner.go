package runner

import "context"

type Runer interface {
	Run(ctx context.Context) (finalTaskUID int, err error)
}
