package hippo

import (
	"fmt"
	"slices"

	"github.com/hkoosha/giraffe/internal/gstrings"
)

type regEntry struct {
	fn      *Fn
	aliases []string
}

func (r *regEntry) clone() regEntry {
	return regEntry{
		fn:      r.fn,
		aliases: slices.Clone(r.aliases),
	}
}

func (r *regEntry) String() string {
	return fmt.Sprintf(
		"FnEntry[fn=%s, aliases=%s]",
		r.fn.typ,
		gstrings.Joined(r.aliases),
	)
}
