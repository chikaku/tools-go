package network

import (
	"bytes"
	crand "crypto/rand"
	"io"
	"math/rand"
	"net"
	"testing"
)

func TestTcpTunnel(t *testing.T) {
	const (
		addrEcho         = "127.0.0.1:23301"
		addrTunnelClient = "127.0.0.1:23302"
		addrTunnelServer = "127.0.0.1:23303"
	)

	{
		// echo server
		lsn, err := net.Listen("tcp", addrEcho)
		if err != nil {
			t.Error(err)
		}

		go func() {
			conn, err := lsn.Accept()
			if err != nil {
				t.Error(err)
			}

			defer conn.Close()
			io.Copy(conn, conn)
		}()
	}

	{
		s, err := NewTCPTunnelServer(addrTunnelServer, addrEcho)
		if err != nil {
			t.Error(err)
		}
		go func() {
			if err := s.Serve(); err != nil {
				t.Error(err)
			}
		}()
	}

	{
		c, err := NewTCPTunnelAgent(addrTunnelClient, addrTunnelServer)
		if err != nil {
			t.Error(err)
		}
		go func() {
			if err := c.Serve(); err != nil {
				t.Error(err)
			}
		}()
	}

	conn, err := net.Dial("tcp", addrTunnelClient)
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 128; i++ {
		n := rand.Intn(1 << 20)
		b0 := make([]byte, n)
		crand.Read(b0)

		_, err = conn.Write(b0)
		if err != nil {
			t.Error(err)
		}

		b1 := make([]byte, n)
		_, err = io.ReadFull(conn, b1)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(b0, b1) {
			t.Error("wrong response")
		}
	}
}
