package gtesting

import (
	"net"
	"strconv"
	"testing"
	"time"

	. "github.com/hkoosha/giraffe/internal/dot0"
)

func NewSocketServer(
	t *testing.T,
	port int,
	handler func(net.Conn, error),
) func() {
	t.Helper()
	EnsureTesting()

	ln, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(port)))
	if err != nil {
		panic(E(err))
	}

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
		if err := ln.Close(); err != nil {
			panic(E(err))
		}
	}
}
