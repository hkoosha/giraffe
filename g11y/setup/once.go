package setup

type Bypassed interface {
	SetBypassed(bool)
}

type Then interface {
	Then(...string) OnceHandle
}

type OnceHandle interface {
	Then

	Finish() OnceHandle
	EnsureOpen() OnceHandle
	EnsureDone() OnceHandle
}

type OnceRegistry interface {
	Bypassed
	Then
}

func NewOnceRegistry() OnceRegistry {
	return newOnceRegistry()
}

func Global() OnceRegistry {
	return global
}

// ============================================================================.

func SetBypassed(bypassed bool) {
	global.SetBypassed(bypassed)
}

func Finish(
	what ...string,
) OnceHandle {
	return global.then(what).Finish()
}

func EnsureOpen(
	what ...string,
) OnceHandle {
	return global.then(what).EnsureOpen()
}

func EnsureDone(
	what ...string,
) OnceHandle {
	return global.then(what).EnsureDone()
}
