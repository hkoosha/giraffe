package toggles

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/t11y/glog"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/toggles/internal"
)

type Toggler interface {
	internal.Sealed

	Query(
		gtx.Context,
		string,
		...Value,
	) (bool, error)
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
	Value() any

	Str() (string, bool)
	Bln() (bool, bool)
	I64() (int64, bool)
	U64() (uint64, bool)
}

type Storage interface {
	Get(
		gtx.Context,
		string,
		Values,
	) (*bool, error)
}

type Values []Value

func (v Values) Assoc() map[string]any {
	m := make(map[string]any, len(v))
	for _, val := range v {
		m[val.Name()] = val.Value()
	}
	return m
}

// ====================================.

func Router(
	defaultCase Condition,
	togglers ...Storage,
) Toggler {
	return newRouter(defaultCase, togglers)
}

func Ephemeral(
	lg glog.Lg,
) *InMemory {
	return newInMemory(lg)
}

func Constant(
	lg glog.Lg,
	enabled bool,
) Storage {
	return newConstant(lg, enabled)
}

// ====================================.

//goland:noinspection GoUnusedExportedFunction
func Always() Condition {
	return yes_
}

//goland:noinspection GoUnusedExportedFunction
func Never() Condition {
	return no_
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
