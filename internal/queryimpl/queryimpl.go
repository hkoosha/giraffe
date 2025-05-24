package queryimpl

import (
	"fmt"

	"github.com/hkoosha/giraffe/cmd"
	"github.com/hkoosha/giraffe/dialects"
)

// MaxDepth must fit in the gqflag.QFlag in the sequence part, i.e., 8 bits.
const MaxDepth = 255

type QueryImpl interface {
	fmt.Stringer

	// Resolved(func(query string) (data string, _ error)) (QueryImpl, error)

	Flags() cmd.QFlag
	Dialect() dialects.Dialect
	Escaped() string

	Attr() string
	Index() int

	Root() QueryImpl
	Leaf() QueryImpl
	Prev() QueryImpl
	Next() QueryImpl

	WithoutOverwrite() QueryImpl
	WithOverwrite() QueryImpl
	WithMake() QueryImpl
}
