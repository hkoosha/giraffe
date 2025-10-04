package dialects

import (
	"errors"

	"github.com/hkoosha/giraffe/qcmd"
)

const (
	Giraffe1v1 Dialect = "giraffe1v1"
	Http1v1    Dialect = "http1v1"

	Unknown Dialect = ""
)

var errUnknown = errors.New("dialect unknown")
var errMismatch = errors.New("dialect mismatch")

type Dialect string

func (d Dialect) String() string {
	return string(d)
}

func (d Dialect) matches(
	spec string,
) (
	explicit bool,
	matched bool,
) {
	switch {
	case len(spec) < 1 || spec[0] != qcmd.Dialect.Byte():
		matched = d == Giraffe1v1
		explicit = false

	default:
		to := len(d) + 1
		matched = len(spec) >= to && spec[1:to] == string(d)
		explicit = matched
	}

	return explicit, matched
}

func dialectOf(
	spec string,
) (_ Dialect, explicit bool, _ error) {
	if ex, ok := Giraffe1v1.matches(spec); ok {
		return Giraffe1v1, ex, nil
	}
	if ex, ok := Http1v1.matches(spec); ok {
		return Http1v1, ex, nil
	}

	return Unknown, false, errUnknown
}

func DialectOf(
	spec string,
) (Dialect, error) {
	d, _, err := dialectOf(spec)
	return d, err
}

func Normalized(
	spec string,
) (Dialect, string, error) {
	d, explicit, err := dialectOf(spec)
	if err != nil {
		return Unknown, "", err
	}

	if !explicit {
		spec = qcmd.Dialect.String() + d.String() + qcmd.Sep.String() + spec
	}

	return d, spec, nil
}

func ErrUnknown() error {
	return errUnknown
}

func ErrMismatch() error {
	return errMismatch
}
