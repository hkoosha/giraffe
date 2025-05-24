package internal

import (
	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

type Seal struct{}

type Sealed interface {
	private() Seal
}

type Sealer struct{}

func (q *Sealer) private() Seal {
	panic(EF("do not call"))
}
