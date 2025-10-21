package gtx

import (
	"slices"
	"sync"
)

type events struct {
	store []any
	mu    *sync.Mutex
}

func (e *events) get() []any {
	e.mu.Lock()
	defer e.mu.Unlock()

	return slices.Clone(e.store)
}

func (e *events) add(v any) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.store = append(e.store, v)
}
