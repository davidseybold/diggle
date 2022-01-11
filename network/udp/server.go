package udp

import (
	"fmt"
	"net"

	"github.com/davidseybold/dns-resolver/resolver"
)

type UDPServer struct {
	r *resolver.Resolver
}

func (u *UDPServer) Listen() error {

	sAddr, err := net.ResolveUDPAddr("udp4", ":53")
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", sAddr)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)

	for {
		n, cAddr, err := conn.ReadFromUDP(buffer)

		res := u.handleRequest(buffer[:n-1])

		conn.WriteToUDP(res, cAddr)

		if err != nil {
			fmt.Println("error occurred", err)
		}
	}
}

func (u *UDPServer) handleRequest(buffer []byte) []byte {
	return nil
}
