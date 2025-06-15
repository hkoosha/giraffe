package ztime

import (
	"time"

	"github.com/hkoosha/giraffe/zebra/ztime/internal"
)

const (
	LayoutRFC3339NanoNoTz = "2006-01-02T15:04:05.999999999"
	LayoutRFC3339NoTz     = "2006-01-02T15:04:05"
)

// ==============================================================================.

type RFC3339NoTz struct {
	time.Time
}

//goland:noinspection GoMixedReceiverTypes
func (t RFC3339NoTz) MarshalJSON() ([]byte, error) {
	return internal.Marshal(LayoutRFC3339NoTz, t.Time)
}

//goland:noinspection GoMixedReceiverTypes
func (t *RFC3339NoTz) UnmarshalJSON(b []byte) error {
	parsed, err := internal.Unmarshal(LayoutRFC3339NoTz, b)
	if err != nil {
		return err
	}

	*t = RFC3339NoTz{parsed}

	return nil
}

// ==============================================================================.

type RFC3339 struct {
	time.Time
}

//goland:noinspection GoMixedReceiverTypes
func (t RFC3339) MarshalJSON() ([]byte, error) {
	return internal.Marshal(time.RFC3339, t.Time)
}

//goland:noinspection GoMixedReceiverTypes
func (t *RFC3339) UnmarshalJSON(b []byte) error {
	parsed, err := internal.Unmarshal(time.RFC3339, b)
	if err != nil {
		return err
	}

	*t = RFC3339{parsed}

	return nil
}

// ==============================================================================.

type RFC3339Nano struct {
	time.Time
}

//goland:noinspection GoMixedReceiverTypes
func (t RFC3339Nano) MarshalJSON() ([]byte, error) {
	return internal.Marshal(time.RFC3339Nano, t.Time)
}

//goland:noinspection GoMixedReceiverTypes
func (t *RFC3339Nano) UnmarshalJSON(b []byte) error {
	parsed, err := internal.Unmarshal(time.RFC3339Nano, b)
	if err != nil {
		return err
	}

	*t = RFC3339Nano{parsed}

	return nil
}

// ==============================================================================.

type RFC3339NanoNoTz struct {
	time.Time
}

//goland:noinspection GoMixedReceiverTypes
func (t RFC3339NanoNoTz) MarshalJSON() ([]byte, error) {
	return internal.Marshal(LayoutRFC3339NanoNoTz, t.Time)
}

//goland:noinspection GoMixedReceiverTypes
func (t *RFC3339NanoNoTz) UnmarshalJSON(b []byte) error {
	parsed, err := internal.Unmarshal(LayoutRFC3339NanoNoTz, b)
	if err != nil {
		return err
	}

	*t = RFC3339NanoNoTz{parsed}

	return nil
}
