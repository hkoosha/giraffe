package setup

import (
	"sync"
)

type contextKeyT int

var keyOnce = New()
var keyMu = sync.Mutex{}
var contextKey contextKeyT = 1111743090

func Key(
	name string,
) any {
	keyOnce.Finish("giraffe", "core", "setup", "context_key", name)

	keyMu.Lock()
	defer keyMu.Unlock()

	contextKey = contextKey + 1

	return contextKey
}
