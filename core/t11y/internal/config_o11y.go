package internal

import (
	"os"
	"regexp"
	"sync/atomic"

	"github.com/hkoosha/giraffe/core/t11y/env"
)

var (
	IsDebug       = atomic.Bool{}
	IsToString    = atomic.Bool{}
	IsTracer      = atomic.Bool{}
	IsUnsafeError = atomic.Bool{}

	envFlagRe = regexp.MustCompile(`^\s*(?i)on|yes|true|enabled|en|1\s*$`)
)

func init() {
	_ = IsDebug.CompareAndSwap(false, chk(env.Debug))
	_ = IsToString.CompareAndSwap(false, chk(env.ToString))
	_ = IsTracer.CompareAndSwap(false, chk(env.Tracer))
	_ = IsUnsafeError.CompareAndSwap(false, chk(env.UnsafeErrors))
}

func chk(envVar string) bool {
	e := os.Getenv(envVar)

	return envFlagRe.MatchString(e)
}
