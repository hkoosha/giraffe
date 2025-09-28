package headers_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	_ "github.com/hkoosha/giraffe/conn/headers"
	"github.com/hkoosha/giraffe/helper/testhelper/enumtelper"
)

func TestHeaders(t *testing.T) {
	t.Skip("casing func in enum helper pkg is not working yet")

	t.Run("content_types casing", func(t *testing.T) {
		pkg, err := packages.Load(&enumtelper.ReadPkgCfg, "../")
		require.NoError(t, err)

		enums, err := enumtelper.Extract(pkg)
		require.NoError(t, err)

		err = enumtelper.CheckWith(
			enums,
			enumtelper.NoIgnore,
			enumtelper.NoOverwrite,
			func(v string) string {
				fixed := enumtelper.DashedTitleCasing(v)
				fixed = http.CanonicalHeaderKey(fixed)
				return fixed
			},
		)
		require.NoError(t, err)
	})
}
