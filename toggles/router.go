package toggles

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hkoosha/giraffe/glog"
)

type Strategy int

func (s Strategy) IsStop(
	_ error,
) bool {
	switch s {
	case Continue:
		return false

	case BailOut:
		return true

	default:
		panic("unknown strategy")
	}
}

const (
	Continue Strategy = iota
	BailOut
)

// ============================================================================.

func NewRouter(
	lg glog.Lg,
	t ...Toggler,
) *Router {
	return &Router{
		lg:       lg,
		togglers: t,
	}
}

var _ Toggler = (*Router)(nil)

type Router struct {
	lg       glog.Lg
	togglers []Toggler

	onQueryErr Strategy
}

func (r *Router) shallow() *Router {
	return *&r
}

func (r *Router) WithOnQueryErr(
	s Strategy,
) *Router {
	cp := r.shallow()
	cp.onQueryErr = s
	return cp
}

// ============================================================================.

func (r *Router) Get(
	ctx context.Context,
	name string,
	attrs ...Attr,
) (bool, error) {
	var err error

	for _, t := range r.togglers {
		var en bool
		if en, err = t.Get(ctx, name, attrs...); err != nil {
			r.lg.Error("failed to query toggle",
				glog.Of("name", name),
				glog.Of("attrs", attrs),
				err,
			)
			if r.onQueryErr.IsStop(err) {
				break
			}
		} else if en {
			return true, nil
		}
	}

	return false, err
}

func (r *Router) GetOrFalse(
	ctx context.Context,
	name string,
	attrs ...Attr,
) bool {
	en, err := r.Get(ctx, name, attrs...)
	return err == nil && en
}

func (r *Router) Enable(
	ctx context.Context,
	name string,
	attrs ...Attr,
) error {
	return r.Set(ctx, name, true, attrs...)
}

func (r *Router) Disable(
	ctx context.Context,
	name string,
	attrs ...Attr,
) error {
	return r.Set(ctx, name, false, attrs...)
}

func (r *Router) Set(
	ctx context.Context,
	name string,
	enabled bool,
	attrs ...Attr,
) error {
	var eg errgroup.Group
	for _, t := range r.togglers {
		eg.Go(func() error {
			err := t.Set(ctx, name, enabled, attrs...)
			if err != nil {
				r.lg.Error("failed to enable",
					glog.Of("name", name),
					glog.Of("attrs", attrs),
					err,
				)
			}
			return err
		})
	}

	return eg.Wait()
}
