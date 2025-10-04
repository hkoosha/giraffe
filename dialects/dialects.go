package dialects

import (
	"errors"

	"github.com/hkoosha/giraffe/qcmd"
	. "github.com/hkoosha/giraffe/t11y/dot"
)

const (
	Giraffe1v1 Dialect = "giraffe1v1"

	Unknown Dialect = ""
)

var (
	errUnknown  = errors.New("dialect unknown")
	errMismatch = errors.New("dialect mismatch")
)

type Dialect string

func (d Dialect) String() string {
	return string(d)
}

//nolint:nonamedreturns
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

//nolint:nonamedreturns
func dialectOf(
	spec string,
) (_ Dialect, explicit bool, _ error) {
	if ex, ok := Giraffe1v1.matches(spec); ok {
		return Giraffe1v1, ex, nil
	}

	return Unknown, false, ErrUnknown()
}

func DialectOf(
	spec string,
) (Dialect, error) {
	d, _, err := dialectOf(spec)
	return d, err
}

//nolint:nonamedreturns
func Denormalized(
	spec string,
) (_ Dialect, denormalizedSpec string, _ error) {
	dialect, explicit, err := dialectOf(spec)
	if err != nil {
		return Unknown, "", err
	}

	if !explicit {
		spec = qcmd.Dialect.String() + dialect.String() + qcmd.Sep.String() + spec
	}

	return dialect, spec, nil
}

//nolint:nonamedreturns
func Normalized(
	spec string,
) (_ Dialect, normalizedSpec string, _ error) {
	dialect, explicit, err := dialectOf(spec)
	if err != nil {
		return Unknown, "", err
	}

	if explicit {
		spec = spec[len(qcmd.Dialect.String())+len(dialect.String())+1:]
	}

	return dialect, spec, nil
}

func ErrUnknown() error {
	return E(errUnknown)
}

func ErrMismatch() error {
	return E(errMismatch)
}
