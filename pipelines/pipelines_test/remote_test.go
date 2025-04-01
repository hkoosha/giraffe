package pipelines_test

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/g11y"
	"github.com/hkoosha/giraffe/pipelines"
)

func mkEkran() pipelines.Server {
	g11y.EnableTracer()
	g11y.EnableUnsafeError()

	mkStep := func(
		step int,
	) pipelines.Fn {
		return func(
			_ context.Context,
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
		}
	}

	reg := pipelines.FnRegistry{}.
		MustWithNamed("fn0", mkStep(0)).
		MustWithNamed("fn1", mkStep(1)).
		MustWithNamed("fn2", mkStep(2))

	ekran := pipelines.NewServer(reg, map[string]pipelines.Plan{
		"plan0": pipelines.Plan{}.
			MergeRegistry(reg).
			MustWithNextNamed("fn0").
			MustWithNextNamed("fn1").
			MustWithNextNamed("fn2"),
	})

	return ekran
}

func TestServer_Ekran(t *testing.T) {
	g11y.EnableTracer()
	g11y.EnableUnsafeError()

	t.Run("ekran", func(t *testing.T) {
		ekran := mkEkran()

		out := bytes.Buffer{}
		req := pipelines.Request{
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

	remote := func(
		ctx context.Context,
		dat giraffe.Datum,
	) (giraffe.Datum, error) {
		u64, err := dat.QU64("meow")
		if err != nil {
			return giraffe.OfErr(), err
		}

		return giraffe.Of1(Q("meow2"), u64*2), nil
	}

	local := func(
		context.Context,
		giraffe.Datum,
	) (giraffe.Datum, error) {
		return giraffe.Of1(Q("fn0"), 111), nil
	}

	mkSrv := func() *httptest.Server {
		reg := pipelines.FnRegistry{}.MustWithNamed("thingy", remote)

		pSrv := pipelines.NewServer(reg, map[string]pipelines.Plan{
			"thingy": pipelines.Plan{}.
				MergeRegistry(reg).
				MustWithNextNamed("thingy"),
		})

		srv := httptest.NewServer(pSrv)

		return srv
	}

	mkCln := func(
		srv *httptest.Server,
	) *pipelines.RunnerFn {
		plan := pipelines.Plan{}.
			WithNext("fn0", local).
			WithNext("rm", pipelines.Remote(
				srv.URL,
				"thingy",
				srv.Client(),
			))

		return M(pipelines.Runner(plan))
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
