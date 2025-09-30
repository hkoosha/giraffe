package glog

const (
	TypAny = iota + 111
	TypStr = iota + 111
	TypInt = iota + 111
	TypErr = iota + 111
)

type Field struct {
	Name  string
	Value any
	Typ   int
}

// =============================================================================

func OfN(name string, v any) Field {
	switch vt := v.(type) {
	case string:
		return OfStringN(name, vt)

	case int:
		return OfIntN(name, vt)

	case error:
		return OfErrN(name, vt)

	default:
		return Field{
			Value: v,
			Typ:   TypAny,
		}
	}
}

func OfIntN(name string, v int) Field {
	return Field{
		Name:  name,
		Value: v,
		Typ:   TypInt,
	}
}

func OfErrN(name string, v error) Field {
	return Field{
		Name:  name,
		Value: v,
		Typ:   TypErr,
	}
}

func OfStringN(name string, v string) Field {
	return Field{
		Name:  name,
		Value: v,
		Typ:   TypStr,
	}
}

// =============================================================================

func Of(v any) Field {
	return OfN("", v)
}

func OfInt(v int) Field {
	return OfN("", v)
}

func OfErr(v error) Field {
	return OfErrN("", v)
}

func OfString(v string) Field {
	return OfStringN("", v)
}
