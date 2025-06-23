package internal

type Seal struct{}

type Sealed interface {
	seal() Seal
}

type Sealer struct{}

func (q *Sealer) seal() Seal {
	panic("do not call")
}
