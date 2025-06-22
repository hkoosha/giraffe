package toggles

import (
	"context"

	"golang.org/x/sync/errgroup"
)

var _ Toggler = (*router)(nil)

type router struct {
	togglers []Toggler
}

func (r *router) Query(
	ctx context.Context,
	name string,
	values ...Value,
) (bool, error) {
	var err error

	for _, t := range r.togglers {
		var en bool
		if en, err = t.Query(ctx, name, values...); err == nil && en {
			return true, nil
		}
	}

	return false, err
}

func (r *router) Enable(
	ctx context.Context,
	name string,
	rest ...Condition,
) error {
	return r.Set(ctx, name, true, rest...)
}

func (r *router) Disable(
	ctx context.Context,
	name string,
	rest ...Condition,
) error {
	return r.Set(ctx, name, false, rest...)
}

func (r *router) Set(
	ctx context.Context,
	name string,
	enable bool,
	rest ...Condition,
) error {
	var eg errgroup.Group
	for _, t := range r.togglers {
		eg.Go(func() error {
			return t.Set(ctx, name, enable, rest...)
		})
	}

	return eg.Wait()
}
