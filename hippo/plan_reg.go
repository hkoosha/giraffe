package hippo

import (
	"maps"
	"slices"
	"strings"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	"github.com/hkoosha/giraffe/hippo/internal/privnames"
	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/g"
	"github.com/hkoosha/giraffe/typing"
	"github.com/hkoosha/giraffe/zebra/z"
)

//goland:noinspection GoUnusedGlobalVariable
var FnRegistry_ = &FnRegistry{
	scope:  nil,
	byType: make(map[typing.Type]regEntry),
}

type FnRegistry struct {
	scope  *giraffe.Query
	byType map[typing.Type]regEntry
}

func (r FnRegistry) String() string {
	const prefix = "FnRegistry["
	const suffix = "]"

	value := strings.Builder{}

	value.WriteString("scope=")
	if r.scope != nil {
		value.WriteString(r.scope.String())
	} else {
		value.WriteString("nil")
	}

	names := z.ItFlatten(z.Apply2AsV(r.byType, func(it regEntry) []string {
		return it.aliases
	}))

	value.WriteString(", names=[")
	value.WriteString(g.JoinIt(names))
	value.WriteString("]")

	value.WriteString(", types=[")
	value.WriteString(g.Join(z.ItApplied(maps.Keys(r.byType), func(it typing.Type) string {
		return it.String()
	})...))
	value.WriteString("]")

	return prefix + value.String() + suffix
}

func (r FnRegistry) clone() FnRegistry {
	cp := FnRegistry{
		scope:  r.scope,
		byType: maps.Clone(r.byType),
	}

	if cp.byType == nil {
		cp.byType = make(map[typing.Type]regEntry)
	}

	return cp
}

func (r FnRegistry) WithNamed(
	name string,
	fn *Fn_,
) (FnRegistry, error) {
	if !fn.IsValid() {
		panic(EF("invalid fn"))
	}

	if !privnames.SimpleName.MatchString(name) {
		return fnRegistryErr, hippoerr.NewPlanInvalidFnName(name)
	}

	for _, f := range r.byType {
		for _, a := range f.aliases {
			if a == name {
				return fnRegistryErr, hippoerr.NewRegDuplicateFnError(name)
			}
		}
	}

	cp := r.clone()

	entry, ok := cp.byType[fn.typ]

	if ok {
		entry = entry.clone()
		entry.aliases = append(entry.aliases, name)
	} else {
		entry = regEntry{
			fn:      fn,
			aliases: []string{name},
		}
	}

	cp.byType[fn.typ] = entry

	return cp, nil
}

func (r FnRegistry) hasNamed(
	alias string,
) bool {
	for _, e := range r.byType {
		if slices.Contains(e.aliases, alias) {
			return true
		}
	}

	return false
}

func (r FnRegistry) has(
	fn *Fn_,
	alias string,
) (bool, error) {
	if fn == nil {
		panic(EF("nil fn"))
	}

	existing, ok := r.byType[fn.typ]

	switch {
	case !ok:
		return false, nil

	case slices.Contains(existing.aliases, alias):
		return true, nil

	case r.hasNamed(alias):
		return false, hippoerr.NewRegDuplicateFnError(alias)

	default:
		return false, nil
	}
}

func (r FnRegistry) With(
	fn *Fn_,
	aliases ...string,
) (FnRegistry, error) {
	if fn == nil {
		panic(EF("nil fn"))
	}

	if _, ok := r.byType[fn.typ]; ok {
		return fnRegistryErr, hippoerr.NewRegDuplicateFnError(fn.typ.String())
	}

	cp := r.clone()
	cp.byType[fn.typ] = regEntry{
		fn:      fn,
		aliases: aliases,
	}

	return cp, nil
}

func (r FnRegistry) MustWithNamed(
	name string,
	fn *Fn_,
) FnRegistry {
	return M(r.WithNamed(name, fn))
}

func (r FnRegistry) MustWith(
	fn *Fn_,
	aliases ...string,
) FnRegistry {
	return M(r.With(fn, aliases...))
}

func (r FnRegistry) Named(
	name string,
) (*Fn_, error) {
	for _, e := range r.byType {
		if slices.Contains(e.aliases, name) {
			return e.fn, nil
		}
	}

	return nil, hippoerr.NewRegMissingFnError(typing.OfErr(), name)
}

func (r FnRegistry) Typed(
	ty typing.Type,
) (*Fn_, error) {
	if r.byType == nil {
		return nil, hippoerr.NewRegMissingFnError(ty, "")
	}

	if e, ok := r.byType[ty]; ok {
		return e.fn, nil
	}

	return nil, hippoerr.NewRegMissingFnError(ty, "")
}

func (r FnRegistry) Merge(
	other FnRegistry,
) (FnRegistry, error) {
	cp := r.clone()

	for k, oFn := range other.byType {
		entry, ok := cp.byType[k]

		switch {
		case ok && entry.fn.typ != oFn.fn.typ:
			return fnRegistryErr, EF(
				"cannot merge due to fn mismatch: %s != %s",
				entry.String(),
				oFn.String(),
			)

		case ok:
			entry = entry.clone()
			entry.aliases = append(entry.aliases, oFn.aliases...)

		default:
			entry = oFn
		}

		cp.byType[k] = entry
	}

	return cp, nil
}

func (r FnRegistry) Dump() FnRegistry {
	return r
}
