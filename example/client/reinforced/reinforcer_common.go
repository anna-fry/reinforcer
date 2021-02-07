// Code generated by reinforcer, DO NOT EDIT.

package reinforced

import (
	"context"
	goresilience "github.com/slok/goresilience"
)

type base struct {
	errorPredicate func(string, error) bool
	runnerFactory  runnerFactory
}
type runnerFactory interface {
	GetRunner(name string) goresilience.Runner
}

var RetryAllErrors = func(_ string, _ error) bool {
	return true
}

type Option func(*base)

func WithRetryableErrorPredicate(fn func(string, error) bool) Option {
	return func(o *base) {
		o.errorPredicate = fn
	}
}
func (b *base) run(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	return b.runnerFactory.GetRunner(name).Run(ctx, fn)
}
