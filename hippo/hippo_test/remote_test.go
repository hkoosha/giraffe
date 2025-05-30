package hippo_test

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"
	. "github.com/hkoosha/giraffe/internal/dot0"
	. "github.com/hkoosha/giraffe/internal/dot1"
)

func mkEkran(
	t *testing.T,
) remote.Server {
	t.Helper()

	g11y.EnableTracer()
	g11y.EnableUnsafeError()

	mkStep := func(
		step int,
	) *hippo.Fn_ {
		fn := hippo.MustFnOf0(func(
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

	reg := hippo.FnRegistry{}.
		MustWithNamed("fn0", mkStep(0)).
		MustWithNamed("fn1", mkStep(1)).
		MustWithNamed("fn2", mkStep(2))

	ekran, err := remote.NewServer(reg, map[string]*hippo.Plan{
		"plan0": hippo.Plan_.
			MustAndRegistry(reg).
			MustWithNextNamed("fn0").
			MustWithNextNamed("fn1").
			MustWithNextNamed("fn2"),
	})

	require.NoError(t, err)

	return ekran
}

func TestServer_Ekran(t *testing.T) {
	g11y.EnableTracer()
	g11y.EnableUnsafeError()

	t.Run("ekran", func(t *testing.T) {
		ekran := mkEkran(t)

		out := bytes.Buffer{}
		req := remote.Request{
			Compensations: nil,
			Init:          map[string]any{"m-1": 123},
			Plan:          "plan0",
		}

		err := ekran(t.Context(), bytes.NewReader(M(json.Marshal(req))), &out)
		require.NoError(t, err, "ekran failed: %s", g11y.FmtStacktraceOf(err))

		var fin any
		OK(json.Unmarshal(out.Bytes(), &fin))

		t.Log(M(giraffe.Make(fin)).Pretty())
	})
}

func TestServer_Http(t *testing.T) {
	g11y.EnableTracer()
	g11y.EnableUnsafeError()

	remoteFn := hippo.MustFnOf0(func(
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		u64, err := dat.QU64("meow")
		if err != nil {
			return giraffe.OfErr(), err
		}

		return giraffe.Of1(Q("meow2"), u64*2), nil
	})

	local := hippo.MustFnOf0(func(
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return giraffe.Of1(Q("fn0"), 111), nil
	})

	mkSrv := func() *httptest.Server {
		reg := hippo.FnRegistry{}.MustWithNamed("thingy", remoteFn)

		pSrv, err := remote.NewServer(reg, map[string]*hippo.Plan{
			"thingy": hippo.Plan_.
				MustAndRegistry(reg).
				MustWithNextNamed("thingy"),
		})

		require.NoError(t, err)

		srv := httptest.NewServer(pSrv)

		return srv
	}

	mkCln := func(
		srv *httptest.Server,
	) *hippo.PipelineFn {
		plan := hippo.Plan_.
			MustWithNext("fn0", local).
			MustWithNext("rm", remote.Remote(
				srv.URL,
				"thingy",
				srv.Client(),
			))

		return M(hippo.Pipeline(plan))
	}

	t.Run("server", func(t *testing.T) {
		srv := mkSrv()
		defer srv.Close()

		cln := mkCln(srv)

		fin, err := cln.Ekran(t.Context(), giraffe.Of1("meow", 333))
		require.NoError(t, err)

		t.Log(fin.Pretty())
	})
}
