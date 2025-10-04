package queryimpl

import (
	"fmt"

	"github.com/hkoosha/giraffe/qflag"
)

// MaxDepth must fit in the gqflag.QFlag in the sequence part, i.e., 8 bits.
const MaxDepth = 255

type QueryImpl interface {
	fmt.Stringer

	Flags() qflag.QFlag

	Plus(query string) (QueryImpl, error)

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
