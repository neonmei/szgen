package runner

import (
	"context"
)

type Task interface {
	Execute(context.Context) error
	Name() string
}
