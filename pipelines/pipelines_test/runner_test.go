package pipelines_test

import (
	"context"
	"errors"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/pipelines"
)

func fn0(
	context.Context,
	giraffe.Datum,
) (giraffe.Datum, error) {
	return giraffe.Of1(Q("ns0.my_out_fn0"), []uint64{
		11,
		22,
		33,
		44,
		55,
	}), nil
}

func fn1(
	context.Context,
	giraffe.Datum,
) (giraffe.Datum, error) {
	return giraffe.Of1(Q("ns1.my_out_fn1"), []int{2, 4}), nil
}

func fn2(
	_ context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	fn0Out, err := dat.QU64s("ns0.my_out_fn0")
	if err != nil {
		return giraffe.OfErr(), err
	}

	fn1Out, err := dat.QISzs("ns1.my_out_fn1")
	if err != nil {
		return giraffe.OfErr(), err
	}

	sum := uint64(0)
	for _, i := range fn1Out {
		sum += fn0Out[i]
	}

	return giraffe.Of1(
		Q("sum"),
		sum,
	), nil
}

func mul(
	step int,
) pipelines.Fn {
	return func(
		_ context.Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		out := "m" + strconv.Itoa(step)
		in := "m"
		if step > 0 {
			in += strconv.Itoa(step - 1)
		}

		sum, err := dat.QInt(Q(in))
		if err != nil {
			return giraffe.OfErr(), err
		}

		sum.Mul(sum, big.NewInt(3))

		return giraffe.Of1(
			Q(out),
			sum,
		), nil
	}
}

func fail(
	msg string,
) pipelines.Fn {
	return func(
		context.Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return giraffe.OfErr(), errors.New("I failed like this: " + msg)
	}
}

func write(
	t *testing.T,
	path string,
	content string,
) {
	t.Helper()

	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0o644,
	)
	require.NoError(t, err)

	defer file.Close()
	_, err = file.WriteString(content)
	require.NoError(t, err)
}

func TestRunner(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		plan := pipelines.Plan{}.
			WithNext("my_fn0", fn0).
			WithNext("my_fn1", fn1).
			WithNext("my_fn2", fn2)

		pipeline, err := pipelines.Runner(plan)
		require.NoError(t, err)

		state := M(pipeline.Ekran(t.Context(), giraffe.OfEmpty()))

		fin, err := state.Get("fin.sum")
		require.NoError(t, err)

		finCast, err := fin.U64()
		require.NoError(t, err)
		assert.Equal(t, uint64(88), finCast)

		t.Logf("final state:\n%s\n", state.Pretty())

		write(
			t,
			"/tmp/giraffe_state.inspection.json",
			state.Pretty()+"\n",
		)

		write(
			t,
			"/tmp/giraffe_fin.inspection.json",
			fin.Pretty()+"\n",
		)
	})
}

func TestRunner_Compensation(t *testing.T) {
	t.Run("compensation by error message", func(t *testing.T) {
		plan := pipelines.Plan{}.
			WithNext("f_0", fail("thingy")).
			WithNext("m_1", mul(1)).
			WithCompensator(
				pipelines.Compensator{}.
					ForErrorWith(
						regexp.MustCompile("thingy"),
						giraffe.Of1("m0", 101),
					),
			)

		pipeline := M(pipelines.Runner(plan))
		state := M(pipeline.Ekran(t.Context(), giraffe.Of1("m", 33)))
		fin := M(state.QU64("fin.m1"))

		assert.Equal(t, uint64(303), fin)

		t.Logf("final state:\n%s\n", state.Pretty())

		write(
			t,
			"/tmp/giraffe_state.inspection.json",
			state.Pretty()+"\n",
		)

		write(
			t,
			"/tmp/giraffe_fin.inspection.json",
			state.Pretty()+"\n",
		)
	})

	t.Run("compensation by step", func(t *testing.T) {
		plan := pipelines.Plan{}.
			WithNext("m_0", mul(0)).
			WithNext("m_1", mul(1)).
			WithNext("m_2", mul(2)).
			WithNext("f_0", fail("thingy")).
			WithNext("m_4", mul(4)).
			WithCompensator(
				pipelines.Compensator{}.
					ForStepWith(
						3,
						giraffe.Of1("m3", 101),
					),
			)

		pipeline := M(pipelines.Runner(plan))
		state := M(pipeline.Ekran(t.Context(), giraffe.Of1("m", 33)))
		fin := M(state.QU64("fin.m4"))

		assert.Equal(t, uint64(303), fin)

		t.Logf("final state:\n%s\n", state.Pretty())

		write(
			t,
			"/tmp/giraffe_state.inspection.json",
			state.Pretty()+"\n",
		)

		write(
			t,
			"/tmp/giraffe_fin.inspection.json",
			state.Pretty()+"\n",
		)
	})
}
