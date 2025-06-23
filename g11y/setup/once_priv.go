package setup

import (
	"runtime/debug"
	"strings"
	"sync"

	. "github.com/hkoosha/giraffe/internal/dot0"
)

const (
	locked   state = true
	unlocked state = false
)

var bypass = false

var setBypass = sync.Once{}

type state bool

func (s state) String() string {
	if s == locked {
		return "locked"
	}
	return "unlocked"
}

type onceLocks struct {
	sections map[string]string
	mu       sync.Mutex
}

var locks = onceLocks{
	sections: make(map[string]string),
	mu:       sync.Mutex{},
}

func get(key string) (string, state) {
	trace, s := locks.sections[key]

	if state(s) == locked {
		return trace, locked
	}

	return trace, unlocked
}

func lock(key string) {
	locks.sections[key] = string(debug.Stack())
}

func check(
	key string,
	require state,
) {
	if bypass {
		return
	}

	trace, actual := get(key)
	if actual != require {
		const sep = "\n\n===============>"
		panic(EF(
			"once check failed: key=%s, expecting=%s, got=%s;%s\npreviously at:\n%s%s\nnow at:\n%s",
			key,
			require,
			actual,
			sep,
			trace,
			sep,
			string(debug.Stack()),
		))
	}
}

func toKey(
	domain string,
	what ...string,
) string {
	return strings.Join(append([]string{domain}, what...), ".")
}
