package network

import (
	"errors"
	"net"
	"time"

	"github.com/pion/stun"
)

var stunPublicServerList = []string{
	"stun1.l.google.com:19302",
	"142.250.21.127:19302",
	"172.217.213.127:19305",
}

// GetUDPMapAddress return the public UDP address mapped by the NAT.
// The input local UDP address is optinal, and if not provided, a random
// address will be used
func GetUDPMapAddress(local *net.UDPAddr) (*net.UDPAddr, error) {
	for _, addr := range stunPublicServerList {
		remote, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			continue
		}

		openAddr, err := stunGetOpenAddress(local, remote)
		if err == nil {
			return openAddr, nil
		}
	}

	return nil, errors.New("no available STUN service")
}

func stunGetOpenAddress(local, remote *net.UDPAddr) (*net.UDPAddr, error) {
	conn, err := net.DialUDP("udp", local, remote)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	cli, err := stun.NewClient(conn, stun.WithTimeoutRate(time.Second*2))
	if err != nil {
		return nil, err
	}

	var (
		innerErr error
		udpAddr  net.UDPAddr
	)
	if err = cli.Do(
		stun.MustBuild(stun.TransactionID, stun.BindingRequest),
		func(ev stun.Event) {
			if ev.Error != nil {
				innerErr = err
				return
			}

			var xorAddr stun.XORMappedAddress
			if err1 := xorAddr.GetFrom(ev.Message); err1 != nil {
				innerErr = err1
				return
			}

			udpAddr.IP = xorAddr.IP
			udpAddr.Port = xorAddr.Port
		},
	); err != nil {
		return nil, err
	}

	if innerErr != nil {
		return nil, innerErr
	}

	return &udpAddr, nil
}
