package toggles

import (
	"context"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/hkoosha/giraffe/toggles/internal"
)

func newRouter(
	defaultCase Condition,
	togglers []Storage,
) Toggler {
	return &router{
		Sealer: internal.Sealer{},

		defaultCase: defaultCase,
		togglers:    togglers,
	}
}

type router struct {
	internal.Sealer
	defaultCase Condition
	togglers    []Storage
}

func (r *router) Query(
	ctx context.Context,
	name string,
	values ...Value,
) (bool, error) {
	var err error

	for _, t := range r.togglers {
		var en *bool
		if en, err = t.Get(ctx, name, slices.Clone(values)); err == nil && en != nil && *en {
			return true, nil
		}
	}

	return r.defaultCase.test(values), err
}

func (r *router) Enable(
	ctx context.Context,
	name string,
	c Condition,
) error {
	return r.Set(ctx, name, true, c)
}

func (r *router) Disable(
	ctx context.Context,
	name string,
	c Condition,
) error {
	return r.Set(ctx, name, false, c)
}

func (r *router) Set(
	ctx context.Context,
	name string,
	enable bool,
	c Condition,
) error {
	var eg errgroup.Group
	for _, t := range r.togglers {
		eg.Go(func() error {
			return t.Set(ctx, name, enable, c)
		})
	}

	return eg.Wait()
}
