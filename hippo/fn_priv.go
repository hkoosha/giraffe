package hippo

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo/internal/hippoerr"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

var errInvalidFn = errors.New("invalid fn")

func (f *Fn_) ensure() {
	if !f.IsValid() {
		panic(EF("invalid fn"))
	}
}

// =====================================.

func chkDatPresent(
	dat giraffe.Datum,
	keys []giraffe.Query,
) error {
	if len(keys) == 0 {
		return nil
	}

	var missing []giraffe.Query
	for _, k := range keys {
		if !dat.Has(k) {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		return hippoerr.NewFnMissingKeysError(missing...)
	}

	return nil
}

// =====================================.

func (f *Fn_) clone() *Fn_ {
	f.ensure()

	if f == nil {
		return nil
	}

	return &Fn_{
		exe:           f.exe,
		scopedOut:     f.scopedOut,
		scopedIn:      f.scopedIn,
		noOverwriting: f.noOverwriting,
		inputs:        slices.Clone(f.inputs),
		optionals:     slices.Clone(f.optionals),
		outputs:       slices.Clone(f.outputs),
		replicated:    maps.Clone(f.replicated),
		swapped:       maps.Clone(f.swapped),
		selected:      slices.Clone(f.selected),
		typ:           f.typ.Clone(),
		name:          f.name,
	}
}

func (f *Fn_) String() string {
	return fmt.Sprintf("Fn[%s][%s]", f.typ, f.name)
}
