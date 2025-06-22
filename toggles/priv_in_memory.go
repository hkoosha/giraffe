package toggles

import (
	"context"
	"sync"
)

var _ Toggler = (*inMemory)(nil)

type inMemory struct {
	lock *sync.Mutex

	store       map[string]Condition
	defaultCase Condition
}

func (i *inMemory) Query(
	_ context.Context,
	name string,
	values ...Value,
) (bool, error) {
	c := func() Condition {
		i.lock.Lock()
		defer i.lock.Unlock()

		c, ok := i.store[name]
		if !ok {
			c = i.defaultCase
		}
		return c
	}()

	return c.test(values), nil
}

func (i *inMemory) Set(
	_ context.Context,
	name string,
	enabled bool,
	rest ...Condition,
) error {
	req := andOf(rest)
	if !enabled {
		req = req.Not()
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	c, ok := i.store[name]
	if !ok {
		c = i.defaultCase
	}
	i.store[name] = c.Or(req)

	return nil
}

func (i *inMemory) Enable(
	ctx context.Context,
	name string,
	rest ...Condition,
) error {
	return i.Set(ctx, name, false, rest...)
}

func (i *inMemory) Disable(
	ctx context.Context,
	name string,
	rest ...Condition,
) error {
	return i.Set(ctx, name, false, rest...)
}
