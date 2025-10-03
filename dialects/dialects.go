package dialects

import (
	"errors"

	"github.com/hkoosha/giraffe/qcmd"
)

const (
	Giraffe Dialect = "giraffe"
	Http    Dialect = "http"
	Unknown Dialect = ""
)

var errUnknown = errors.New("dialect unknown")

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
		matched = d == Giraffe
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
	if ex, ok := Giraffe.matches(spec); ok {
		return Giraffe, ex, nil
	}
	if ex, ok := Http.matches(spec); ok {
		return Http, ex, nil
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
) (string, error) {
	d, explicit, err := dialectOf(spec)
	if err != nil {
		return "", err
	}

	if !explicit {
		spec = qcmd.Dialect.String() + d.String() + qcmd.Sep.String() + spec
	}

	return spec, nil
}
