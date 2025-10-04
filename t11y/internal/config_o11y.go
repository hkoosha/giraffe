package internal

import (
	"os"
	"regexp"
	"sync/atomic"

	"github.com/hkoosha/giraffe/genv"
)

var (
	IsDebug       = atomic.Bool{}
	IsToString    = atomic.Bool{}
	IsTracer      = atomic.Bool{}
	IsUnsafeError = atomic.Bool{}

	envFlagRe = regexp.MustCompile(`^\s*(?i)on|yes|true|enabled|en|1\s*$`)
)

func init() {
	_ = IsDebug.CompareAndSwap(false, chk(genv.Debug))
	_ = IsToString.CompareAndSwap(false, chk(genv.ToString))
	_ = IsTracer.CompareAndSwap(false, chk(genv.Tracer))
	_ = IsUnsafeError.CompareAndSwap(false, chk(genv.UnsafeErrors))
}

func chk(envVar string) bool {
	e := os.Getenv(envVar)

	return envFlagRe.MatchString(e)
}
