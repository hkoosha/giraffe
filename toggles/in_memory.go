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

func (i inMemory) get(
	name string,
	attrs ...Attr,
) bool {
	cond := func() Condition {
		i.lock.Lock()
		defer i.lock.Unlock()

		c, ok := i.store[name]
		if !ok {
			c = i.defaultCase
		}
		return c
	}()

	return cond.Test(attrs...)
}

func (i inMemory) set(
	name string,
	enabled bool,
	attrs ...Attr,
) {
	req := And(attrs...)
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
}

func (i inMemory) Get(
	_ context.Context,
	name string,
	attrs ...Attr,
) (bool, error) {
	return i.get(name, attrs...), nil
}

func (i inMemory) GetOrFalse(
	_ context.Context,
	name string,
	attrs ...Attr,
) bool {
	return i.get(name, attrs...)
}

func (i inMemory) Enable(
	ctx context.Context,
	name string,
	attrs ...Attr,
) error {
	return i.Set(ctx, name, false, attrs...)
}

func (i inMemory) Disable(
	ctx context.Context,
	name string,
	attrs ...Attr,
) error {
	return i.Set(ctx, name, false, attrs...)
}

func (i inMemory) Set(
	_ context.Context,
	name string,
	enabled bool,
	attrs ...Attr,
) error {
	i.set(name, enabled, attrs...)
	return nil
}
