package toggles

import (
	"context"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/toggles/internal"
)

type Toggler interface {
	Query(
		context.Context,
		string,
		...Value,
	) (bool, error)

	Set(
		context.Context,
		string,
		bool,
		...Condition,
	) error

	Enable(
		context.Context,
		string,
		...Condition,
	) error

	Disable(
		context.Context,
		string,
		...Condition,
	) error
}

type Condition interface {
	condition

	Name() string
	And(...Condition) Condition
	Or(...Condition) Condition
	Not() Condition
}

type Value interface {
	internal.Sealed

	Name() string
	value() any
}

//goland:noinspection GoUnusedExportedFunction
func Always() Condition {
	return yes_
}

func Eq[V giraffe.Basic](
	name string,
	v V,
) Condition {
	return eqOf(name, v)
}

func Gt[V giraffe.Num](
	name string,
	v V,
) Condition {
	return gtOf(name, v)
}

func Lt[V giraffe.Num](
	name string,
	v V,
) Condition {
	return ltOf(name, v)
}

func In[V giraffe.Ord](
	name string,
	v ...V,
) Condition {
	return inOf(name, v)
}

func Of[V giraffe.Basic](
	name string,
	v V,
) Value {
	return valueOf(name, v)
}

func Router(
	t ...Toggler,
) Toggler {
	return &router{
		togglers: t,
	}
}
