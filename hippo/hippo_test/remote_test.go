package hippo_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/contrib/gtestinghippo"
	"github.com/hkoosha/giraffe/core/gtesting"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/hippo"
)

func ekran(
	t *testing.T,
	dat giraffe.Datum,
) giraffe.Datum {
	t.Helper()

	mkStep := func(
		step int,
	) *hippo.Fn {
		fn := hippo.FnOf(func(
			_ gtx.Context,
			call hippo.Call,
		) (giraffe.Datum, error) {
			in := "m" + strconv.Itoa(step-1)
			out := "m" + strconv.Itoa(step)

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
		require.True(t, fn.IsValid())
		return fn
	}

	reg := hippo.MkFnRegistry().
		MustWithNamed("fn0", mkStep(0)).
		MustWithNamed("fn1", mkStep(1)).
		MustWithNamed("fn2", mkStep(2))

	return gtestinghippo.EkranRemote(t, reg, dat, "fn0", "fn1", "fn2")
}

func TestServer_Ekran(t *testing.T) {
	// TODO un-skip
	t.Skip()

	t.Run("ekran", func(t *testing.T) {
		gtesting.Preamble(t)

		fin := ekran(t, Dat(map[string]int{
			"m-1": 123,
		}))

		t.Log(fin.Pretty())
	})
}

func TestServer_Http(t *testing.T) {
	fn := func(
		_ gtx.Context,
		call hippo.Call,
	) (giraffe.Datum, error) {
		dat := call.Data()

		u64, err := dat.QU64("meow")
		if err != nil {
			return giraffe.OfErr(), err
		}
		return giraffe.Of1(Q("meow2"), u64*2), nil
	}

	t.Run("server", func(t *testing.T) {
		gtesting.Preamble(t)

		dat := giraffe.Of1("meow", 333)

		srv := gtestinghippo.MakeTestServer(t, fn)
		defer srv.Close()

		fin := gtestinghippo.Call(t, srv, dat)

		gtesting.Write(t, "fin.json", fin.Pretty())
	})
}
