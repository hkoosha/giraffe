package hippo

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/cmd"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

func simple(dat giraffe.Datum) func(query giraffe.Query) (string, error) {
	return func(query giraffe.Query) (string, error) {
		ss, err := dat.Get(query)
		if err != nil {
			return "", E(
				err,
				EF("missing path variable in context: %s", query.String()),
			)
		}

		sss, err := ss.SimpleString()
		if err != nil {
			return "", E(
				err,
				EF("path variable cannot be formatted to string: %s", query.String()),
			)
		}

		return sss, nil
	}
}

func MkHttpCallFn(
	channel *DatumChannel,
) HttpCallFn {
	return HttpCallFn{
		DatumChannel: channel,
	}
}

type HttpCallFn struct {
	*DatumChannel
}

func (e *HttpCallFn) Fn() *Fn {
	fn := FnOf(e.exe).
		WithOutput(
			Q(HttpOutputBody),
			Q(HttpOutputHeaders),
		)

	for d := range e.headers {
		fn = fn.WithInput(Q(HttpInputHeader + cmd.Sep.String() + d))
	}

	if e.hasBody {
		fn = fn.WithInput(Q(HttpInputBody))
	}

	for _, d := range e.dcPath.pathPartsVar {
		fn = fn.WithInput(d)
	}

	return fn
}

func (e *HttpCallFn) mkPath(
	dat giraffe.Datum,
) ([]string, error) {
	if e.dcPath.isZero() {
		return []string{e.path}, nil
	}

	return e.dcPath.mkPath(simple(dat))
}

func (e *HttpCallFn) getBody(
	dat giraffe.Datum,
) (giraffe.Datum, int, error) {
	if !e.hasBody {
		return OfErr(), 0, nil
	}

	body, err := dat.Get(Q(HttpInputBody))
	if err != nil {
		return OfErr(), 0, err
	}

	l, err := body.Len()
	if err != nil {
		return OfErr(), 0, err
	}

	return body, l, nil
}

func (e *HttpCallFn) exe(
	ctx Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	path, err := e.mkPath(dat)
	if err != nil {
		return OfErr(), err
	}

	withHeaders, err := e.mkHeaders(simple(dat))
	if err != nil {
		return OfErr(), err
	}

	body, l, err := e.getBody(dat)
	if err != nil {
		return OfErr(), err
	}

	cnx := e.cnx.Cfg().AndHeaders(withHeaders).Datum()

	bodyR := &body
	if l == 0 {
		bodyR = nil
	}

	headers, rx, err := cnx.HCall(ctx, bodyR, path...)
	if err != nil {
		return OfErr(), err
	}

	ret := map[string]any{
		HttpOutputHeaders: headers,
		HttpOutputBody:    rx,
	}

	return giraffe.FromJsonable(ret)
}
