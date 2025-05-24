package hippo_test

import (
	_ "embed"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/core/gtesting"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/hippo"
)

//go:embed test_runner.simple.json
var testRunnerSimple string

func mul(step int) *hippo.Fn {
	return hippo.FnOf(func(
		_ gtx.Context,
		call hippo.Call,
	) (giraffe.Datum, error) {
		dat := call.Data()

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
	})
}

func alwaysFail(msg string) *hippo.Fn {
	return hippo.FnOf(func(
		gtx.Context,
		hippo.Call,
	) (giraffe.Datum, error) {
		return giraffe.OfErr(), errors.New("I failed like this: " + msg)
	})
}

func TestRunner(t *testing.T) {
	// TODO un-skip
	t.Skip()

	t.Run("simple", func(t *testing.T) {
		gtesting.Preamble(t)

		plan := hippo.
			MkPlan().
			MustWithNext("my_fn0", hippo.FnOf(func(
				gtx.Context,
				hippo.Call,
			) (giraffe.Datum, error) {
				return giraffe.Of1(
					Q("ns0.my_out_fn0"),
					[]uint64{11, 22, 33, 44, 55},
				), nil
			})).
			MustWithNext("my_fn1", hippo.FnOf(func(
				gtx.Context,
				hippo.Call,
			) (giraffe.Datum, error) {
				return giraffe.Of1(Q("ns1.my_out_fn1"), []int{2, 4}), nil
			})).
			MustWithNext("my_fn2", hippo.FnOf(func(
				_ gtx.Context,
				call hippo.Call,
			) (giraffe.Datum, error) {
				dat := call.Data()

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
			}))

		pipeline, err := hippo.MkPipeline(plan)
		require.NoError(t, err)

		state, err := pipeline.Ekran(gtx.Of(t.Context()), giraffe.OfEmpty())
		require.NoError(t, err)
		gtesting.Write(t, "state.json", state.Pretty())

		fin, err := state.Get("fin.sum")
		require.NoError(t, err)
		gtesting.Write(t, "fin.json", fin.Pretty())

		finCast, err := fin.U64()
		require.NoError(t, err)
		assert.Equal(t, uint64(88), finCast)

		assert.Equal(
			t,
			strings.TrimSpace(testRunnerSimple),
			strings.TrimSpace(state.Pretty()),
		)
	})

	t.Run("compensation by error message", func(t *testing.T) {
		gtesting.Preamble(t)

		plan := hippo.
			MkPlan().
			MustWithNext("f_0", alwaysFail("thingy")).
			MustWithNext("m_1", mul(1)).
			AndCompensator(
				hippo.Compensator{}.
					ForErrorWith(
						regexp.MustCompile("thingy"),
						giraffe.Of1("m0", 101),
					),
			)

		pipeline, err := hippo.MkPipeline(plan)
		require.NoError(t, err)

		state, err := pipeline.Ekran(gtx.Of(t.Context()), giraffe.Of1("m", 33))
		require.NoError(t, err)

		fin, err := state.QU64("fin.m1")
		require.NoError(t, err)

		assert.Equal(t, uint64(303), fin)

		gtesting.Write(t, "state.json", state.Pretty())
		gtesting.Write(t, "fin.json", state.Pretty())
	})

	t.Run("compensation by step", func(t *testing.T) {
		gtesting.Preamble(t)

		plan := hippo.
			MkPlan().
			MustWithNext("m_0", mul(0)).
			MustWithNext("m_1", mul(1)).
			MustWithNext("m_2", mul(2)).
			MustWithNext("f_0", alwaysFail("thingy")).
			MustWithNext("m_4", mul(4)).
			AndCompensator(
				hippo.Compensator{}.
					ForStepWith(
						3,
						giraffe.Of1("m3", 101),
					),
			)

		pipeline, err := hippo.MkPipeline(plan)
		require.NoError(t, err)

		state, err := pipeline.Ekran(gtx.Of(t.Context()), giraffe.Of1("m", 33))
		require.NoError(t, err)

		fin, err := state.QU64("fin.m4")
		require.NoError(t, err)

		assert.Equal(t, uint64(303), fin)

		gtesting.Write(t, "state.json", state.Pretty())
		gtesting.Write(t, "fin.json", state.Pretty())
	})
}
