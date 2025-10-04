package hippoerr

import (
	"github.com/hkoosha/giraffe"
	"github.com/hkoosha/giraffe/internal/g"
)

type FnMissingKeysErrorState struct {
	missing []giraffe.Query
}

func (e *FnMissingKeysErrorState) String(_ *HippoError) string {
	strs := make([]string, len(e.missing))
	for i, m := range e.missing {
		strs[i] = m.String()
	}

	return "missing keys: " + g.Joined(strs)
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
