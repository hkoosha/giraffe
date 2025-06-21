package toggles

import (
	"github.com/hkoosha/giraffe/zebra/z"
)

type Condition interface {
	Test(...Attr) bool

	And(...Condition) Condition
	Or(...Condition) Condition
	Not() Condition

	uid() string
	seal() seal

	isYes() bool
	isNo() bool
}

type Attr interface {
	Condition

	Name() string
}

func And(attrs ...Attr) Condition {
	return andOf(z.Applied(attrs, func(it Attr) Condition {
		return it
	}))
}

func Yes() Condition {
	return &yes{}
}

func No() Condition {
	return &no{}
}
