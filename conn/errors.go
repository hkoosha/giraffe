package conn

import (
	"strconv"
)

//goland:noinspection GoUnusedConst
const (
	ReasonUnexpectedStatusCode FailureReason = 2
	ReasonEmptyResponse        FailureReason = 3
)

type FailureReason uint

type FailedResponseError struct {
	Resp   any
	Reason FailureReason
}

func (e *FailedResponseError) Error() string {
	return "http request failed: " + strconv.FormatUint(uint64(e.Reason), 10)
}

type MissingEndpointError struct {
	Endpoint string
}

func (m *MissingEndpointError) Error() string {
	return "missing endpoint, name=" + m.Endpoint
}
