package gtesting

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func NewSocketServer(
	t *testing.T,
	port int,
	handler func(net.Conn, error),
) func() {
	t.Helper()
	EnsureTesting()

	ln, err := (&net.ListenConfig{
		Control:   nil,
		KeepAlive: -1,
		KeepAliveConfig: net.KeepAliveConfig{
			Enable:   false,
			Idle:     -1,
			Interval: -1,
			Count:    -1,
		},
	}).Listen(t.Context(), "", strconv.Itoa(port))
	require.NoError(t, err)

	go func() {
		conn, lErr := ln.Accept()
		if lErr == nil {
			//goland:noinspection GoUnhandledErrorResult
			defer conn.Close()
		}

		handler(conn, lErr)
	}()

	time.Sleep(10 * time.Millisecond)

	return func() {
		require.NoError(t, ln.Close())
	}
}
