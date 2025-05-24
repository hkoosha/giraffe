package hippo_test

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/typing"
)

var _ FnTyp = (*hippo.Fn)(nil)

type FnOrig = *hippo.Fn

// =============================================================================

type FnTyp interface {
	FnConfigured
	FnNamed

	Type() typing.Type
	IsValid() bool
	String() string

	Dump() FnOrig
}

type FnNamed interface {
	// Name() FnOrig
	Named(string) FnOrig
}

type FnConfiguredExec interface {
	WithoutSkipOnExists() FnOrig
	WithSkipOnExists() FnOrig
	SetSkipOnExists(bool) FnOrig

	Skipped() bool
	WithoutSkipped() FnOrig
	WithSkipped() FnOrig
	SetSkipped(bool) FnOrig

	SkippedWith() (giraffe.Datum, bool)
	WithSkippedWith(d giraffe.Datum) FnOrig
	WithoutSkippedWith() FnOrig
}

type FnConfiguredIn interface {
	AndInputs(...giraffe.Query) FnOrig
	WithInputs(...giraffe.Query) FnOrig
	WithoutInputs() FnOrig

	AndOptionals(...giraffe.Query) FnOrig
	WithOptional(...giraffe.Query) FnOrig
	WithoutOptional() FnOrig
}

type FnConfiguredOut interface {
	AndOutputs(...giraffe.Query) FnOrig
	WithOutput(...giraffe.Query) FnOrig
	WithoutOutput() FnOrig

	WithScope(giraffe.Query) FnOrig
	WithoutScope() FnOrig
}

type FnConfiguredData interface {
	AndCopied(map[giraffe.Query]giraffe.Query) FnOrig
	WithCopied(map[giraffe.Query]giraffe.Query) FnOrig
	WithoutCopied() FnOrig

	WithCombine(map[giraffe.Query][]giraffe.Query) FnOrig
	WithoutCombine() FnOrig

	AndSelect(...giraffe.Query) FnOrig
	WithSelect(...giraffe.Query) FnOrig
	WithoutSelect() FnOrig
}

type FnConfigured interface {
	FnConfiguredExec
	FnConfiguredIn
	FnConfiguredOut
	FnConfiguredData

	/*
			AndSwapping(map[giraffe.Query]giraffe.Query) FnOrig
			WithSwapping(map[giraffe.Query]giraffe.Query) FnOrig
			WithoutSwapping() FnOrig

		    AndArgs(...giraffe.Query) FnOrig
			WithArgs(...giraffe.Query) FnOrig
			WithoutArgs() FnOrig
	*/
}
