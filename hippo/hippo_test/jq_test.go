package hippo

import (
	"testing"

	"github.com/hkoosha/giraffe/contrib/gtestinghippo"
	. "github.com/hkoosha/giraffe/dot"
	"github.com/hkoosha/giraffe/gtesting"
	"github.com/hkoosha/giraffe/hippo"
)

func TestJq(t *testing.T) {
	t.Run("jq", func(t *testing.T) {
		gtesting.Preamble(t)

		fin := gtestinghippo.EkranFn(t,
			Of(map[string]string{
				"foo": "999",
			}),
			M(hippo.MkJqFn(".fooz = 123")).Fn(),
		)

		gtesting.Write(t, "fin.json", fin.Pretty())
	})

}
