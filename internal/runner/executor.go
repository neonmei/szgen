package runner

import (
	"context"
)

type Executor interface {
	Execute(ctx context.Context, tasks []Task) error
}
