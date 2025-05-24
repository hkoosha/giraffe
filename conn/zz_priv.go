package conn

import (
	"io"
	"regexp"
	"slices"
)

var (
	endpointNameRe = regexp.MustCompile(
		"^[a-zA-Z0-9-_]+$")

	endpointAddrRe = regexp.MustCompile(
		`^(http|https)://(?P<addr>[a-zA-Z0-9-_.]{1,255})(:(?P<port>\d{1,5}))?$`)

	endpointAddrReNames = slices.DeleteFunc(
		endpointAddrRe.SubexpNames()[1:],
		func(it string) bool { return it == "" },
	)
)

type noBodyT struct{}

func (noBodyT) Read([]byte) (int, error)         { return 0, io.EOF }
func (noBodyT) Close() error                     { return nil }
func (noBodyT) WriteTo(io.Writer) (int64, error) { return 0, nil }

var nobody = noBodyT{}
