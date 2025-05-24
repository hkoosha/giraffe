package gtestinghippo

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/core/gtesting"
	"github.com/hkoosha/giraffe/core/t11y/gtx"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"

	. "github.com/hkoosha/giraffe/dot"
)

func MakeServer(
	t *testing.T,
	handler *hippo.Fn,
) remote.Server {
	t.Helper()

	reg := hippo.
		MkFnRegistry().
		MustWithNamed("thingy", handler)

	srv, err := remote.NewServer(reg, map[string]*hippo.Plan{
		"thingy": hippo.
			MkPlan().
			MustAndRegistry(reg).
			MustWithNextNamed("thingy"),
	})
	require.NoError(t, err)

	return srv
}

func MakeTestServer(
	t *testing.T,
	exe hippo.Exe,
) *httptest.Server {
	t.Helper()

	srv := MakeServer(t, hippo.FnOf(exe))
	return httptest.NewServer(srv)
}

func Call(
	t *testing.T,
	srv *httptest.Server,
	dat giraffe.Datum,
) giraffe.Datum {
	t.Helper()

	plan := hippo.
		MkPlan().
		MustWithNext("local", hippo.FnOf(func(
			gtx.Context,
			hippo.Call,
		) (giraffe.Datum, error) {
			return giraffe.Of1(Q("fn0"), 111), nil
		})).
		MustWithNext("remote", remote.Remote(
			"thingy",
			conn.MakeCfg(gtesting.Zap(t)).
				WithTransport(srv.Client().Transport).
				AndEndpoint("thingy", srv.URL).
				WithMustEndpointNamed("thingy").
				Raw(),
		))

	cl, err := hippo.MkPipeline(plan)
	gtesting.NoError(t, err)

	fin, err := cl.Ekran(
		gtx.Of(t.Context()),
		dat,
	)
	gtesting.NoError(t, err)

	return fin
}
