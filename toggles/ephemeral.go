package toggles

import (
	"context"
	"sync"

	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/glog"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

var _ Storage = (*InMemory)(nil)

func newInMemory(lg glog.Lg) *InMemory {
	g11y.NonNil(lg)
	_ = lg

	return &InMemory{
		lock:  &sync.Mutex{},
		store: make(map[string]Condition),
	}
}

type InMemory struct {
	lock  *sync.Mutex
	store map[string]Condition
}

func (i *InMemory) get(
	name string,
	values Values,
) *bool {
	c, ok := func() (Condition, bool) {
		i.lock.Lock()
		defer i.lock.Unlock()

		c, ok := i.store[name]
		return c, ok
	}()

	if !ok {
		return nil
	}

	return Ref(c.test(values))
}

func (i *InMemory) set(
	name string,
	enabled bool,
	req Condition,
) {
	if !enabled {
		req = req.Not()
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	c, ok := i.store[name]
	if !ok {
		c = Always()
	}

	i.store[name] = c.And(req)
}

func (i *InMemory) Get(
	_ context.Context,
	name string,
	values Values,
) (*bool, error) {
	return i.get(name, values), nil
}

func (i *InMemory) Set(
	_ context.Context,
	name string,
	enabled bool,
	req Condition,
) error {
	i.set(name, enabled, req)
	return nil
}

// ============================================================================.

type constant struct {
	enabled bool
}

func newConstant(
	lg glog.Lg,
	enabled bool,
) Storage {
	_ = lg
	return &constant{enabled: enabled}
}

func (i *constant) Get(
	context.Context,
	string,
	Values,
) (*bool, error) {
	// Make a copy first.
	en := i.enabled
	return &en, nil
}
