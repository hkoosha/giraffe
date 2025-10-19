package headers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	_ "github.com/hkoosha/giraffe/conn/headers"
	"github.com/hkoosha/giraffe/core/gtesting"
)

func TestHeaders(t *testing.T) {
	t.Skip("casing func in enum helper pkg is not working yet")

	t.Run("content_types casing", func(t *testing.T) {
		pkg, err := packages.Load(gtesting.ReadPkgCfg(), "../")
		require.NoError(t, err)

		enums, err := gtesting.Extract(pkg)
		require.NoError(t, err)

		err = gtesting.CheckWith(
			enums,
			gtesting.NoIgnore(),
			gtesting.NoOverwrite(),
			func(v string) string {
				fixed := gtesting.DashedTitleCasing(v)
				fixed = http.CanonicalHeaderKey(fixed)
				return fixed
			},
		)
		require.NoError(t, err)
	})
}
