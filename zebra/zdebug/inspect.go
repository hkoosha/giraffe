package zdebug

import (
	"encoding/json"
	"fmt"

	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/g11y"
	. "github.com/hkoosha/giraffe/internal/dot0"
)

func Inspect(v any) string {
	switch vt := v.(type) {
	case giraffe.Datum:
		return Inspect(M(vt.Raw()))

	default:
		return string(M(json.MarshalIndent(v, "", "   ")))
	}
}

//nolint:forbidigo
func Dump[V any](v V) V {
	Sep()
	fmt.Println(Inspect(v))

	return v
}

func DumpE[V any](v V, err error) V {
	g11y.Ensure(err)

	return Dump(v)
}

//nolint:forbidigo
func Sep() {
	const sep = "================================================================================"

	fmt.Printf("\n\n%s\n\n", sep)
}
