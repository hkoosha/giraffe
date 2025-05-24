package gtesting

import (
	"errors"
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

func ReadPkgCfg() *packages.Config {
	//nolint:exhaustruct
	//goland:noinspection GoUnusedGlobalVariable
	return &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax,
		Fset: token.NewFileSet(),
		Dir:  ".",
	}
}

func NoIgnore() map[string]struct{} {
	return map[string]struct{}{}
}

func NoOverwrite() map[string]string {
	return map[string]string{}
}

// =============================================================================

//nolint:err113
func Extract(
	pkg []*packages.Package,
) (map[string]string, error) {
	var err error
	packages.Visit(pkg, nil, func(it *packages.Package) {
		if err != nil {
			return
		}

		if it.Errors != nil {
			err = it.Errors[0]
			return
		}

		if it.Module != nil && it.Module.Error != nil {
			err = errors.New("mod error: " + it.Module.Error.Err)
			return
		}
	})

	if err != nil {
		return nil, err
	}

	enums := make(map[string]string)
	for _, e := range pkg {
		for _, n := range e.Types.Scope().Names() {
			//nolint:gocritic
			if c, ok := e.Types.Scope().Lookup(n).(*types.Const); ok {
				name := c.Name()
				value := constant.StringVal(c.Val())

				if _, seen := enums[name]; seen {
					return nil, errors.New("duplicate const: " + name)
				}

				enums[name] = value
			}
		}
	}

	return enums, nil
}

//nolint:err113
func Check(
	enums map[string]string,
	transform func(name, value string) string,
) error {
	if len(enums) == 0 {
		return errors.New("no enums")
	}

	var bad []string
	seenValues := make(map[string]struct{})

	for name, value := range enums {
		expecting := transform(name, value)

		if expecting != name {
			bad = append(bad, fmt.Sprintf(
				"%s: expecting=%s, actual=%s",
				name,
				expecting,
				value,
			))
		}

		if _, ok := seenValues[name]; ok {
			bad = append(bad, fmt.Sprintf(
				"%s=%s: duplicate",
				name,
				value,
			))
		}
	}

	if len(bad) == 0 {
		return nil
	}

	return errors.New(strings.Join(bad, "; "))
}

func CheckWith(
	enums map[string]string,
	ignore map[string]struct{},
	overwrite map[string]string,
	transform func(string) string,
) error {
	return Check(enums, func(name, value string) string {
		if _, ok := ignore[name]; ok {
			return value
		}

		if o, ok := overwrite[name]; ok {
			return o
		}

		return transform(name)
	})
}

// =============================================================================

func DashedTitleCasing(
	s string,
) string {
	alphanum := regexp.MustCompile("[^a-zA-Z0-9]")
	num := regexp.MustCompile(`^\d.+`)

	sb := strings.Builder{}
	sb.Grow(len(s) * 2)

	for _, part := range alphanum.Split(s, -1) {
		if len(part) > 1 {
			part = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
			if num.MatchString(s) {
				part = part[:1] + strings.ToUpper(part[1:2]) + part[2:]
			}
		}

		sb.WriteString(part)
	}

	return sb.String()
}
