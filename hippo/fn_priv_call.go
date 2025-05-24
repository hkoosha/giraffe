package hippo

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/hippo/internal"
)

func mkCall(
	name string,
	data giraffe.Datum,
	args *giraffe.Datum,
) Call {
	return fnCall{
		Sealer: internal.Sealer{},
		args:   args,
		data:   data,
		name:   name,
	}
}

// TODO Provide a argOnce flavor, where does not need func to reparse its args.
type fnCall struct {
	internal.Sealer

	args *giraffe.Datum
	data giraffe.Datum

	// Be careful with this! It makes the fn context-sensitive.
	// (context as in programming languages && context-independent languages)
	name string
}

func (c fnCall) clone() fnCall {
	return fnCall{
		Sealer: c.Sealer,
		args:   c.args,
		data:   c.data,
		name:   c.name,
	}
}

func (c fnCall) CheckPresent(
	dat giraffe.Datum,
	queries []giraffe.Query,
) error {
	var missing []giraffe.Query
	for _, q := range queries {
		if ok, err := dat.Has(q); err != nil {
			return err
		} else if !ok {
			missing = append(missing, q)
		}
	}

	if len(missing) != 0 {
		return giraffe.NewMissingKeyError(missing...)
	}

	return nil
}

func (c fnCall) Data() giraffe.Datum {
	return c.data
}

func (c fnCall) AndData(d giraffe.Datum) (Call, error) {
	merged, err := c.data.Merge(d)
	if err != nil {
		return nil, err
	}

	return c.WithData(merged), nil
}

func (c fnCall) WithData(d giraffe.Datum) Call {
	cp := c.clone()
	cp.data = d
	return cp
}

func (c fnCall) Args() giraffe.Datum {
	if c.args == nil {
		return giraffe.OfEmpty()
	}

	return *c.args
}

func (c fnCall) AndArgs(
	args giraffe.Datum,
) (Call, error) {
	if c.args == nil {
		return c.WithArgs(args), nil
	}

	merged, err := c.args.Merge(args)
	if err != nil {
		return nil, err
	}

	return c.WithArgs(merged), nil
}

func (c fnCall) WithArgs(
	args giraffe.Datum,
) Call {
	cp := c.clone()

	if cp.args == nil {
		cp.args = &args
	}

	return cp
}

func (c fnCall) WithoutArgs() Call {
	cp := c.clone()
	cp.args = nil
	return cp
}

func (c fnCall) Name() string {
	return c.name
}

func (c fnCall) WithName(
	name string,
) Call {
	cp := c.clone()
	cp.name = name
	return cp
}
