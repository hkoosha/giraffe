package gtestinghippo

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/gtesting"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"
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
	srv := MakeServer(t, FnOf(exe))
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
			hippo.Context,
			giraffe.Datum,
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

	cl, err := hippo.Pipeline(plan)
	gtesting.NoError(t, err)

	fin, err := cl.Ekran(
		hippo.ContextOf(t.Context()),
		dat,
	)
	gtesting.NoError(t, err)

	return fin
}
