package internal

type Seal struct{}

type Sealed interface {
	private() Seal
}

type Sealer struct{}

func (q *Sealer) private() Seal {
	panic("do not call")
}
