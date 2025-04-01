package headers_test

import (
	"go/constant"
	"go/token"
	"go/types"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	_ "github.com/hkoosha/giraffe/ghttp/headers"
)

//goland:noinspection SpellCheckingInspection
const headersPkg = "../"

func TestHeaders(t *testing.T) {
	t.Run("header name casing", func(t *testing.T) {
		//nolint:exhaustruct
		cfg := &packages.Config{
			Mode: packages.NeedName |
				packages.NeedTypes |
				packages.NeedTypesInfo |
				packages.NeedSyntax,
			Fset: token.NewFileSet(),
			Dir:  ".",
		}

		pkg, err := packages.Load(cfg, headersPkg)
		require.NoError(t, err)

		packages.Visit(pkg, nil, func(it *packages.Package) {
			require.Empty(t, it.Errors)
			if it.Module != nil {
				require.Nil(t, it.Module.Error)
			}
		})

		for _, p := range pkg {
			scope := p.Types.Scope()
			for _, name := range scope.Names() {
				c, ok := scope.Lookup(name).(*types.Const)
				if !ok {
					continue
				}

				hN := c.Name()
				hV := constant.StringVal(c.Val())

				require.NotEmpty(t, hN)
				require.NotEmpty(t, hV)

				assert.Equal(t, http.CanonicalHeaderKey(hV), hV)
				assert.Equal(t, strings.ReplaceAll(hV, "-", ""), hN)
			}
		}
	})
}
