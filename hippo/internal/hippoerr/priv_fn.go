package hippoerr

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/internal/g"
)

type FnMissingKeysErrorState struct {
	missing []giraffe.Query
}

func (e *FnMissingKeysErrorState) String(_ *HippoError) string {
	return "missing keys: " + g.JoinedFn(e.missing, func(it giraffe.Query) string {
		return it.String()
	})
}

// NewFnMissingKeysError Private function, do not call outside hippo package.
func NewFnMissingKeysError(
	missing []giraffe.Query,
) error {
	return NewHippoError(
		ErrCodeMissingKeys,
		"missing keys",
		&FnMissingKeysErrorState{
			missing: missing,
		},
	)
}
