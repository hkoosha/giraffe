package hippo

import (
	"encoding/json"

	"github.com/hkoosha/giraffe"
)

// FnConfig
//
// TODO mix as private field in Fn itself.
//
//nolint:lll
type FnConfig struct {
	Combine      *map[giraffe.Query][]giraffe.Query `json:"combine,omitempty"        yaml:"combine,omitempty"`
	Gather       *map[giraffe.Query]giraffe.Query   `json:"gather,omitempty"         yaml:"gather,omitempty"`
	Copy         *map[giraffe.Query]giraffe.Query   `json:"copy,omitempty"           yaml:"copy,omitempty"`
	Require      *[]giraffe.Query                   `json:"require,omitempty"        yaml:"require,omitempty"`
	Select       *[]giraffe.Query                   `json:"select,omitempty"         yaml:"select,omitempty"`
	Scoped       *giraffe.Query                     `json:"scoped,omitempty"         yaml:"scoped,omitempty"`
	Args         *giraffe.Datum                     `json:"args,omitempty"           yaml:"args,omitempty"`
	SkipOnExists *bool                              `json:"skip_on_exists,omitempty" yaml:"skip_on_exists,omitempty"`
	Skipped      *bool                              `json:"skipped,omitempty"        yaml:"skipped,omitempty"`
	SkippedWith  *giraffe.Datum                     `json:"skipped_with,omitempty" yaml:"skipped_with,omitempty"`
	NoSkipWith   *giraffe.Datum                     `json:"no_skip_with,omitempty"   yaml:"no_skip_with,omitempty"`

	Fn string `json:"fn"                       yaml:"fn"`

	// Swap         *map[giraffe.Query]giraffe.Query `json:"swap,omitempty"           yaml:"swap,omitempty"`
}

func (f *FnConfig) Validate() []error {
	var errs []error

	if f.Combine != nil {
		for k, vs := range *f.Combine {
			for _, v := range vs {
				if _, err := giraffe.GQParse(k.String()); err != nil {
					errs = append(errs, err)
				}
				if _, err := giraffe.GQParse(v.String()); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}

	if f.Gather != nil {
		for k, v := range *f.Gather {
			if _, err := giraffe.GQParse(k.String()); err != nil {
				errs = append(errs, err)
			}
			if _, err := giraffe.GQParse(v.String()); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if f.Copy != nil {
		for k, v := range *f.Copy {
			if _, err := giraffe.GQParse(k.String()); err != nil {
				errs = append(errs, err)
			}
			if _, err := giraffe.GQParse(v.String()); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if f.Require != nil {
		for _, v := range *f.Require {
			if _, err := giraffe.GQParse(v.String()); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if f.Select != nil {
		for _, v := range *f.Select {
			if _, err := giraffe.GQParse(v.String()); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if f.Scoped != nil {
		if _, err := giraffe.GQParse(f.Scoped.String()); err != nil {
			errs = append(errs, err)
		}
	}

	if f.Args != nil {
		b, err := f.Args.MarshalJSON()
		if err != nil {
			errs = append(errs, err)
		} else {
			var dd giraffe.Datum
			err1 := json.Unmarshal(b, &dd)
			if err1 != nil {
				errs = append(errs, err1)
			}
		}
	}

	return nil
}

func (f *FnConfig) Configure(
	fn *Fn,
) (*Fn, error) {
	if f.Require != nil {
		fn = fn.WithInputs(*f.Require...)
	}

	combine := map[giraffe.Query][]giraffe.Query{}
	if f.Combine != nil {
		combine = *f.Combine
	}
	if f.Gather != nil {
		for k, v := range *f.Gather {
			mixed := append(combine[k], v)
			combine[k] = mixed
		}
	}
	if f.Gather != nil || f.Combine != nil {
		fn = fn.WithCombine(combine)
	}

	if f.Select != nil {
		fn = fn.WithSelect(*f.Select...)
	}

	if f.Copy != nil {
		fn = fn.WithCopied(*f.Copy)
	}

	// if f.Swap != nil {
	// 	fn = fn.WithSwapping(*f.Swap)
	// }

	if f.Scoped != nil {
		fn = fn.WithScope(*f.Scoped)
	}

	if f.SkipOnExists != nil {
		fn = fn.SetSkipOnExists(*f.SkipOnExists)
	}

	if f.Skipped != nil {
		fn = fn.SetSkipped(*f.Skipped)
	}

	if f.SkippedWith != nil {
		fn = fn.WithSkippedWith(*f.SkippedWith)
	}

	return fn, nil
}
