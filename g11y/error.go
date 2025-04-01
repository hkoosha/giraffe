package g11y

import (
	"slices"
	"sync"
)

const tracedMsg = "giraffe traced error"

type TracedError struct {
	Mu          *sync.Mutex
	Stacktraces [][]byte
}

func (e *TracedError) Error() string {
	return tracedMsg
}

func (e *TracedError) Get() [][]byte {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	return slices.Clone(e.Stacktraces)
}

func (e *TracedError) Add(stacktrace []byte) {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	e.Stacktraces = append(e.Stacktraces, stacktrace)
}

func NonNil(values ...any) {
	for _, value := range values {
		if value != nil {
			panic("nil value")
		}
	}
}
