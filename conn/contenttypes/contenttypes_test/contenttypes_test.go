package contenttypes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	_ "github.com/hkoosha/giraffe/conn/contenttypes"
	"github.com/hkoosha/giraffe/testhelper"
)

//goland:noinspection SpellCheckingInspection
var exceptions = map[string]string{
	"ApplicationA2l":                          "ApplicationA2L",
	"ApplicationSmpte336m":                    "ApplicationSmpte336M",
	"ApplicationVnd1000mindsDecisionModelXml": "ApplicationVnd1000MindsDecisionModelXml",
	"ApplicationVnd3Gpp2BcmcsinfoXml":         "ApplicationVnd3gpp2BcmcsinfoXml",
}

func TestContentTypes(t *testing.T) {
	t.Skip("casing func in enum helper pkg is not working yet")

	t.Run("content_types casing", func(t *testing.T) {
		pkg, err := packages.Load(testhelper.ReadPkgCfg(), "../")
		require.NoError(t, err)

		enums, err := testhelper.Extract(pkg)
		require.NoError(t, err)

		err = testhelper.CheckWith(
			enums,
			testhelper.NoIgnore(),
			exceptions,
			testhelper.DashedTitleCasing,
		)
		require.NoError(t, err)
	})
}
