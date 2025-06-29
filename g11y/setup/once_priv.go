package setup

import (
	"regexp"
	"runtime/debug"
	"strings"
	"sync"

	. "github.com/hkoosha/giraffe/internal/dot0"
)

const (
	sep      = "::"
	locked   = true
	unlocked = false
)

const (
	stateBypassUnset = iota + 11
	stateNoBypass
	stateBypass
)

var keyRe = regexp.MustCompile("^([a-zA-Z0-9_]+::)*[a-zA-Z0-9_]+$")

var global = newOnceRegistry()

// ============================================================================.

func join(
	what []string,
) string {
	for _, part := range what {
		if !keyRe.MatchString(part) {
			panic(EF("invalid key part, key=%s invalid_part=%s",
				strings.Join(what, sep), part))
		}
	}

	return strings.Join(what, sep)
}

func joinWith(
	key string,
	what []string,
) string {
	return join(append(strings.Split(key, sep), what...))
}

// ============================================================================.

func newOnceRegistry() *onceRegistry {
	return &onceRegistry{
		bypass:    stateBypassUnset,
		traces:    make(map[string]string),
		directory: make(map[string]bool),
		mu:        &sync.RWMutex{},
	}
}

type onceRegistry struct {
	bypass    int
	traces    map[string]string
	directory map[string]bool
	mu        *sync.RWMutex
}

func (o *onceRegistry) _require(
	key string,
	expecting bool,
) {
	if o.bypass == stateBypass || expecting == o.directory[key] {
		return
	}

	here := string(debug.Stack())
	before := o.traces[key]

	expectingStr := "locked"
	if expecting == unlocked {
		expectingStr = "unlocked"
	}

	const lineSep = "\n\n===============>\n"
	panic(EF(
		"once check failed, was expecting: %s==%s"+
			"%s"+
			"previously at:\n%s"+
			"%s"+
			"now at:\n%s",
		key,
		expectingStr,
		lineSep,
		before,
		lineSep,
		here,
	))
}

func (o *onceRegistry) require(
	key string,
	expecting bool,
) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	o._require(key, expecting)
}

func (o *onceRegistry) finish(
	key string,
) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o._require(key, unlocked)

	o.directory[key] = true

	prefix := key + sep
	for it := range o.directory {
		if strings.HasPrefix(it, prefix) {
			o.directory[it] = true
		}
	}
}

func (o *onceRegistry) then(
	what []string,
) handle {
	return handle{
		reg: o,
		key: join(what),
	}
}

func (o *onceRegistry) Then(what ...string) OnceHandle {
	return o.then(what)
}

func (o *onceRegistry) SetBypassed(bypassed bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.bypass > stateBypassUnset {
		isBypassed := true
		if o.bypass == stateNoBypass {
			isBypassed = false
		}
		panic(EF("bypass already set to %t", isBypassed))
	}

	if bypassed {
		o.bypass = stateBypass
	} else {
		o.bypass = stateNoBypass
	}
}

// ====================================.

type handle struct {
	reg *onceRegistry
	key string
}

func (a handle) withChildren(what []string) handle {
	return handle{
		reg: a.reg,
		key: joinWith(a.key, what),
	}
}

func (a handle) Finish() OnceHandle {
	a.reg.finish(a.key)
	return a
}

func (a handle) EnsureOpen() OnceHandle {
	a.reg.require(a.key, unlocked)
	return a
}

func (a handle) EnsureDone() OnceHandle {
	a.reg.require(a.key, locked)
	return a
}

func (a handle) Then(
	what ...string,
) OnceHandle {
	return a.withChildren(what)
}
