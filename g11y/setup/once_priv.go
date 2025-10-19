package setup

import (
	"regexp"
	"runtime/debug"
	"slices"
	"strings"
	"sync"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

const (
	sep      = "::"
	locked   = true
	unlocked = false
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

// ============================================================================.

func newOnceRegistry() *registry {
	return &registry{
		bypass:    nil,
		traces:    make(map[string]string),
		directory: make(map[string]bool),
		mu:        &sync.RWMutex{},
	}
}

var _ Registry = (*registry)(nil)

type registry struct {
	traces    map[string]string
	directory map[string]bool
	mu        *sync.RWMutex
	bypass    *bool
}

func (o *registry) SetBypassed(
	bypassed bool,
) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.bypass != nil {
		panic(EF("bypass already set to %t", *o.bypass))
	}

	o.bypass = &bypassed
}

// ===.

func (o *registry) _requireKey(
	key string,
	expecting bool,
) {
	if o.bypass != nil && *o.bypass == true || expecting == o.directory[key] {
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

func (o *registry) requireKey(
	key string,
	expecting bool,
) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	o._requireKey(key, expecting)
}

func (o *registry) finishKey(
	key string,
) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o._requireKey(key, unlocked)

	o.directory[key] = true

	prefix := key + sep
	for it := range o.directory {
		if strings.HasPrefix(it, prefix) {
			o.directory[it] = true
		}
	}
}

func (o *registry) ensureDoneKey(
	key string,
) {
	o.requireKey(key, locked)
}

func (o *registry) ensureOpenKey(
	key string,
) {
	o.requireKey(key, unlocked)
}

func (o *registry) At(
	what ...string,
) Registry {
	// Validate.
	join(what)

	return &prefixed{
		reg:    o,
		prefix: what,
	}
}

func (o *registry) then(
	what []string,
) handle {
	return handle{
		reg: o,
		key: join(what),
	}
}

func (o *registry) Finish(
	what ...string,
) Handle {
	h := o.then(what)
	h.Finish()
	return h
}

func (o *registry) EnsureOpen(
	what ...string,
) Handle {
	h := o.then(what)
	o.ensureOpenKey(h.key)
	return h
}

func (o *registry) EnsureDone(
	what ...string,
) Handle {
	h := o.then(what)
	o.ensureDoneKey(h.key)
	return h
}

// ====================================.

var _ Registry = (*prefixed)(nil)

type prefixed struct {
	reg    *registry
	prefix []string
}

func (p prefixed) At(what ...string) Registry {
	prefix := append(slices.Clone(p.prefix), what...)
	return p.reg.At(prefix...)
}

func (p prefixed) Finish(what ...string) Handle {
	prefix := append(slices.Clone(p.prefix), what...)
	return p.reg.Finish(prefix...)
}

func (p prefixed) EnsureOpen(what ...string) Handle {
	prefix := append(slices.Clone(p.prefix), what...)
	return p.reg.EnsureOpen(prefix...)
}

func (p prefixed) EnsureDone(what ...string) Handle {
	prefix := append(slices.Clone(p.prefix), what...)
	return p.reg.EnsureDone(prefix...)
}

// ====================================.

type handle struct {
	reg *registry
	key string
}

func (a handle) Finish() {
	a.reg.finishKey(a.key)
}

func (a handle) EnsureOpen() {
	a.reg.ensureOpenKey(a.key)
}

func (a handle) EnsureDone() {
	a.reg.ensureDoneKey(a.key)
}

func (a handle) At(
	what ...string,
) Handle {
	prefix := join(append(strings.Split(a.key, sep), what...))
	return handle{
		reg: a.reg,
		key: prefix,
	}
}
