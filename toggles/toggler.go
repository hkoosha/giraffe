package toggles

import (
	"context"
	"fmt"

	"github.com/hkoosha/giraffe"
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
	sealed

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
	return &eq{
		val: v,
		cond: cond{
			name: name,
			uid_: fmt.Sprintf("eq(%s, %v)", name, v),
		},
	}
}

func Gt[V giraffe.Num](
	name string,
	v V,
) Condition {
	return &gt[V]{
		val: v,
		cond: cond{
			name: name,
			uid_: fmt.Sprintf("gt(%s, %v)", name, v),
		},
	}
}

func Lt[V giraffe.Num](
	name string,
	v V,
) Condition {
	return &lt[V]{
		val: v,
		cond: cond{
			name: name,
			uid_: fmt.Sprintf("lt(%s, %v)", name, v),
		},
	}
}

func In[V giraffe.Ord](
	name string,
	v ...V,
) Condition {
	if len(v) == 0 {
		return no_
	} else if len(v) == 1 {
		return Eq(name, v[0])
	} else if len(v) < 8 {
		return &in[V]{
			val: v,
			cond: cond{
				name: name,
				uid_: fmt.Sprintf("in(%s, %s)", name, uidOf(v)),
			},
		}
	} else {
		return &search[V]{
			val: v,
			cond: cond{
				name: name,
				uid_: fmt.Sprintf("search(%s, %s)", name, uidOf(v)),
			},
		}
	}
}

func Of[V giraffe.Basic](
	name string,
	v V,
) Value {
	return &value{
		name: name,
		val:  v,
	}
}

func Router(
	t ...Toggler,
) Toggler {
	return &router{
		togglers: t,
	}
}
