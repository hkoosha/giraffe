package gtestinghippo

import (
	"bytes"
	"testing"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/gtesting"
	"github.com/hkoosha/giraffe/core/serdes/gson"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"

	. "github.com/hkoosha/giraffe/dot"
)

func EkranFn(
	t *testing.T,
	dat giraffe.Datum,
	fn *hippo.Fn,
) giraffe.Datum {
	t.Helper()

	plan := hippo.
		MkPlan().
		MustWithNext(
			"my_exe",
			fn,
		)

	return Ekran(t, plan, dat)
}

func EkranExe(
	t *testing.T,
	dat giraffe.Datum,
	exe hippo.Exe,
) giraffe.Datum {
	t.Helper()

	return EkranFn(t, dat, hippo.FnOf(exe))
}

func Ekran(
	t *testing.T,
	plan *hippo.Plan,
	dat giraffe.Datum,
) giraffe.Datum {
	t.Helper()

	pipeline, err := hippo.MkPipeline(plan)
	gtesting.NoError(t, err)

	state, err := pipeline.Ekran(gtx.Of(t.Context()), dat)
	gtesting.NoError(t, err)

	return state
}

func Ekran0(
	t *testing.T,
	plan *hippo.Plan,
) giraffe.Datum {
	t.Helper()

	return Ekran(t, plan, OfEmpty())
}

func EkranRemote(
	t *testing.T,
	reg *hippo.FnRegistry,
	dat giraffe.Datum,
	fns ...string,
) giraffe.Datum {
	t.Helper()

	templates := map[string]*hippo.Plan{}
	{
		plan := hippo.MkPlan().MustAndRegistry(reg)
		for _, fn := range fns {
			plan = plan.MustWithNextNamed(fn)
		}
		templates["plan0"] = plan
	}

	ekran, err := remote.NewServer(reg, templates)
	gtesting.NoError(t, err)

	out := bytes.Buffer{}
	req := remote.Request{
		Compensations: nil,
		Init:          dat,
		Plan:          "plan0",
	}

	err = ekran(
		gtx.Of(t.Context()),
		bytes.NewReader(gson.MustMarshal(req)),
		&out,
	)
	gtesting.NoError(t, err)

	deser, err := gson.Unmarshal[any](out.Bytes())
	gtesting.NoError(t, err)

	fin, err := giraffe.From(deser)
	gtesting.NoError(t, err)

	return fin
}
