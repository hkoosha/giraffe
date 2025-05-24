package hippo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/conn"
	"github.com/hkoosha/giraffe/contrib/gtestinghippo"
	"github.com/hkoosha/giraffe/core/gtesting"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/hippo"
	"github.com/hkoosha/giraffe/hippo/remote"
)

func TestRunner_Http(t *testing.T) {
	// TODO unskip
	t.Skip()

	t.Run("fetch", func(t *testing.T) {
		gtesting.Preamble(t)

		srv := httptest.NewServer(http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			resp := ` { "m4": 312 } `
			_, err := w.Write([]byte(resp))
			assert.NoError(t, err)
		}))
		defer srv.Close()

		fn := hippo.MkTunnel(
			"foo",
			conn.MakeCfg(gtesting.Zap(t)).
				WithTxSerde(remote.RequestSerde()).
				WithRxSerde(giraffe.DatumSerde()).
				AndEndpoint("local", srv.URL).
				Datum(),
		)

		plan := hippo.
			MkPlan().
			MustWithNext("http_args", M(hippo.StaticOf(
				P("endpoint", "local"),
				P("path", "/"),
				P("headers", map[string]string{
					"dummy": "value",
				}),
				P("query", map[string]string{
					"q0": "qv0",
					"q1": "qv1",
				}),
			))).
			MustWithNext("http_fn_0", fn.Fn())

		state := gtestinghippo.Ekran0(t, plan)

		fin, err := state.QU64("fin.body.m4")
		require.NoError(t, err)

		assert.Equal(t, uint64(312), fin)

		gtesting.Write(t, "state.json", state.Pretty())
		gtesting.Write(t, "fin.json", state.Pretty())
	})
}
