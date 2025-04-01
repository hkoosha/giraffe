package g11y

import (
	"github.com/hkoosha/giraffe/g11y/internal"
)

func IsDebug() bool {
	return internal.IsDebug.Load()
}

func EnableDebug() {
	internal.IsDebug.Store(true)
}

func DisableDebug() {
	internal.IsDebug.Store(false)
}

// =====================================.

func IsDebugToString() bool {
	return IsDebug() && internal.IsToString.Load()
}

func EnableDebugToString() {
	EnableDebug()
	internal.IsToString.Store(true)
}

func DisableDebugToString() {
	internal.IsToString.Store(false)
}

// =====================================.

func IsTracer() bool {
	return IsDebug() && internal.IsTracer.Load()
}

func EnableTracer() {
	EnableDebug()
	internal.IsTracer.Store(true)
}

func DisableTracer() {
	internal.IsTracer.Store(false)
}

// =====================================.

func IsUnsafeError() bool {
	return internal.IsUnsafeError.Load()
}

func EnableUnsafeError() {
	internal.IsUnsafeError.Store(true)
}

func DisableUnsafeError() {
	internal.IsUnsafeError.Store(false)
}
