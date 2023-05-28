package network

import (
	"log"
	"net"

	"github.com/xtaci/smux"
)

type TCPTunnelServer struct {
	lsn   net.Listener
	raddr string
}

func NewTCPTunnelServer(laadr, raddr string) (*TCPTunnelServer, error) {
	lsn, err := net.Listen("tcp", laadr)
	if err != nil {
		return nil, err
	}

	s := &TCPTunnelServer{
		lsn:   lsn,
		raddr: raddr,
	}

	return s, nil
}

func (s *TCPTunnelServer) Shutdown() error {
	return s.lsn.Close()
}

func (s *TCPTunnelServer) Serve() error {
	for {
		conn, err := s.lsn.Accept()
		if err != nil {
			return err
		}

		sess, err := smux.Server(conn, nil)
		if err != nil {
			conn.Close()
			return err
		}

		go func() {
			if err := s.handleSession(sess); err != nil {
				log.Println(err)
			}
		}()
	}
}

func (s *TCPTunnelServer) handleSession(sess *smux.Session) error {
	for {
		stream, err := sess.AcceptStream()
		if err != nil {
			return err
		}

		conn, err := net.Dial("tcp", s.raddr)
		if err != nil {
			stream.Close()
			return err
		}

		go func() {
			if err := copyStream(conn, stream); err != nil {
				log.Println(err)
			}
		}()
	}
}
