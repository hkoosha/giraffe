package contenttypes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	_ "github.com/hkoosha/giraffe/conn/contenttypes"
	"github.com/hkoosha/giraffe/core/gtesting"
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
		pkg, err := packages.Load(gtesting.ReadPkgCfg(), "../")
		require.NoError(t, err)

		enums, err := gtesting.Extract(pkg)
		require.NoError(t, err)

		err = gtesting.CheckWith(
			enums,
			gtesting.NoIgnore(),
			exceptions,
			gtesting.DashedTitleCasing,
		)
		require.NoError(t, err)
	})
}
