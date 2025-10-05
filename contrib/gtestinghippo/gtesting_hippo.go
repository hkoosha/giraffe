package gtestinghippo

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/contrib/gtesting"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"
)

func Ekran(
	t *testing.T,
	plan *hippo.Plan,
	dat giraffe.Datum,
) giraffe.Datum {
	t.Helper()

	pipeline, err := hippo.Pipeline(plan)
	gtesting.NoError(t, err)

	state, err := pipeline.Ekran(hippo.ContextOf(t.Context()), dat)
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
	reg hippo.FnRegistry,
	dat giraffe.Datum,
	planFns ...string,
) giraffe.Datum {
	t.Helper()

	plan := hippo.MkPlan().MustAndRegistry(reg)

	for _, fn := range planFns {
		plan = plan.MustWithNextNamed(fn)
	}

	templates := map[string]*hippo.Plan{
		"plan0": plan,
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
		hippo.ContextOf(t.Context()),
		bytes.NewReader(M(json.Marshal(req))),
		&out,
	)
	gtesting.NoError(t, err)

	var deser any
	err = json.Unmarshal(out.Bytes(), &deser)
	gtesting.NoError(t, err)

	fin, err := giraffe.From(deser)
	gtesting.NoError(t, err)

	return fin
}
