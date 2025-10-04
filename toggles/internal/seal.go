package internal

import (
	. "github.com/hkoosha/giraffe/internal/dot0"
)

type Seal struct{}

type Sealed interface {
	private() Seal
}

type Sealer struct{}

func (q *Sealer) private() Seal {
	panic(EF("do not call"))
}
