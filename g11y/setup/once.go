package setup

type Bypassed interface {
	SetBypassed(bool)
}

type Handle interface {
	At(...string) Handle

	Finish()
	EnsureOpen()
	EnsureDone()
}

type Registry interface {
	At(...string) Registry

	Finish(...string) Handle
	EnsureOpen(...string) Handle
	EnsureDone(...string) Handle
}

func New() Registry {
	return newOnceRegistry()
}

func Bypassable() (Registry, Bypassed) {
	reg := newOnceRegistry()
	return reg, reg
}

func Global() Registry {
	return global
}

// ============================================================================.

func SetBypassed(bypassed bool) {
	global.SetBypassed(bypassed)
}

func Finish(
	what ...string,
) Handle {
	return global.Finish(what...)
}

func EnsureOpen(
	what ...string,
) Handle {
	return global.EnsureOpen(what...)
}

func EnsureDone(
	what ...string,
) Handle {
	return global.EnsureDone(what...)
}

func At(
	what ...string,
) Registry {
	return global.At(what...)
}
