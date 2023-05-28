package network

import (
	"io"
	"log"
	"net"

	"github.com/chikaku/tools-go"
	"github.com/xtaci/smux"
)

type TCPTunnelAgent struct {
	lsn  net.Listener
	sess *smux.Session
}

func NewTCPTunnelAgent(laadr, raddr string) (*TCPTunnelAgent, error) {
	lsn, err := net.Listen("tcp", laadr)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", raddr)
	if err != nil {
		lsn.Close()
		return nil, err
	}

	sess, err := smux.Client(conn, nil)
	if err != nil {
		lsn.Close()
		conn.Close()
		return nil, err
	}

	agent := &TCPTunnelAgent{
		lsn:  lsn,
		sess: sess,
	}

	return agent, nil
}

func (s *TCPTunnelAgent) Shutdown() error {
	s.sess.Close()
	return s.lsn.Close()
}

func (s *TCPTunnelAgent) Serve() error {
	for {
		conn, err := s.lsn.Accept()
		if err != nil {
			return err
		}

		stream, err := s.sess.OpenStream()
		if err != nil {
			return err
		}

		go func() {
			if err := copyStream(conn, stream); err != nil {
				log.Println(err)
			}
		}()
	}
}

func copyStream(conn net.Conn, stream *smux.Stream) error {
	defer conn.Close()
	defer stream.Close()

	quit := make(chan error)

	go func() {
		if _, err := io.Copy(conn, stream); err != nil {
			tools.SendNonBlock(quit, err)
		}
	}()

	go func() {
		if _, err := io.Copy(stream, conn); err != nil {
			tools.SendNonBlock(quit, err)
		}
	}()

	return <-quit
}
