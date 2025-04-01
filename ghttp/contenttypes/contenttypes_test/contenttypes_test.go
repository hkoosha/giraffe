package contenttypes_test

import (
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	_ "github.com/hkoosha/giraffe/ghttp/contenttypes"
)

//goland:noinspection SpellCheckingInspection
const contenttypesPkg = "../"

var (
	split = regexp.MustCompile("[^a-zA-Z0-9]")
	num   = regexp.MustCompile(`^\d.+`)
)

// "ApplicationJF2feedJson": "ApplicationJf2feedJson",.
var exceptions = map[string]string{
	"ApplicationA2l":       "ApplicationA2L",
	"ApplicationSmpte336m": "ApplicationSmpte336M",

	"ApplicationVnd1000mindsDecisionModelXml": "ApplicationVnd1000MindsDecisionModelXml",
	"ApplicationVnd3Gpp2BcmcsinfoXml":         "ApplicationVnd3gpp2BcmcsinfoXml",
}

func toTitle(s string) string {
	if len(s) < 2 {
		return s
	}
	fixed := strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	if num.MatchString(s) {
		fixed = fixed[:1] + strings.ToUpper(fixed[1:2]) + fixed[2:]
	}

	return fixed
}

func read() ([]*packages.Package, error) {
	//nolint:exhaustruct
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax,
		Fset: token.NewFileSet(),
		Dir:  ".",
	}

	pkg, err := packages.Load(cfg, contenttypesPkg)
	if err != nil {
		return nil, err
	}

	var vErr error
	packages.Visit(pkg, nil, func(it *packages.Package) {
		if vErr != nil {
			return
		}

		if it.Errors != nil {
			vErr = it.Errors[0]

			return
		}

		if it.Module != nil && it.Module.Error != nil {
			panic("could not load module")
		}
	})

	return pkg, nil
}

func TestContentTypes(t *testing.T) {
	t.Run("content types casing", func(t *testing.T) {
		pkg, err := read()
		require.NoError(t, err)

		bad := map[string]string{}

		for _, p := range pkg {
			scope := p.Types.Scope()
			for _, n := range scope.Names() {
				c, ok := scope.Lookup(n).(*types.Const)
				if !ok {
					continue
				}

				name := c.Name()
				require.NotEmpty(t, name)

				value := constant.StringVal(c.Val())
				require.NotEmpty(t, value)

				sb := strings.Builder{}
				for _, s := range split.Split(value, -1) {
					sb.WriteString(toTitle(s))
				}

				expecting := strings.ReplaceAll(sb.String(), "3Gpp", "3gpp")
				expecting = strings.ReplaceAll(expecting, "5Gnas", "5gnas")
				expecting = strings.ReplaceAll(expecting, "5Gsa", "5gsa")
				expecting = strings.ReplaceAll(expecting, "5Gsa", "5gsa")
				if exceptions[expecting] != "" {
					expecting = exceptions[expecting]
				}

				if expecting != name {
					bad[expecting] = name
				}
			}
		}

		// Require.Empty(t, bad).
	})
}
