package hippo_test

import (
	"math/big"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/contrib/gtesting"
	"github.com/hkoosha/giraffe/contrib/gtestinghippo"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

func ekran(
	t *testing.T,
	dat any,
) giraffe.Datum {
	t.Helper()

	mkStep := func(
		step int,
	) *hippo.Fn {
		fn := hippo.MustFnOf(func(
			_ hippo.Context,
			dat giraffe.Datum,
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

	return gtestinghippo.EkranRemote(t, reg, dat, "f0", "f1", "f2")
}

func TestServer_Ekran(t *testing.T) {
	t.Run("ekran", func(t *testing.T) {
		gtesting.Preamble(t)

		fin := ekran(t, map[string]any{
			"m-1": 123,
		})

		t.Log(fin.Pretty())
	})
}

func TestServer_Http(t *testing.T) {
	fnOnRemote := hippo.MustFnOf(func(
		_ hippo.Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		u64, err := dat.QU64("meow")
		if err != nil {
			return giraffe.OfErr(), err
		}

		return giraffe.Of1(Q("meow2"), u64*2), nil
	})

	fnOnLocal := hippo.MustFnOf(func(
		hippo.Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return giraffe.Of1(Q("fn0"), 111), nil
	})

	mkSrv := func() *httptest.Server {
		reg := hippo.
			MkFnRegistry().
			MustWithNamed("thingy", fnOnRemote)

		pSrv, err := remote.NewServer(reg, map[string]*hippo.Plan{
			"thingy": hippo.
				MkPlan().
				MustAndRegistry(reg).
				MustWithNextNamed("thingy"),
		})
		require.NoError(t, err)

		return httptest.NewServer(pSrv)
	}

	mkClient := func(
		srv *httptest.Server,
	) (*hippo.PipelineFn, error) {
		plan := hippo.
			MkPlan().
			MustWithNext("fn0", fnOnLocal).
			MustWithNext("rm", remote.Remote(srv.URL, "thingy", srv.Client()))

		return hippo.Pipeline(plan)
	}

	t.Run("server", func(t *testing.T) {
		gtesting.Preamble(t)

		srv := mkSrv()
		defer srv.Close()

		cln, err := mkClient(srv)
		require.NoError(t, err)

		fin, err := cln.Ekran(hippo.ContextOf(t.Context()), giraffe.Of1("meow", 333))
		require.NoError(t, err)
		gtesting.Write(t, "fin.json", fin.Pretty())
	})
}
