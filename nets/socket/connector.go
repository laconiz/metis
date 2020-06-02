package socket

import (
	"github.com/laconiz/metis/nets/connector"
	"github.com/laconiz/metis/nets/session"
	"net"
)

func Connector() connector.Dialer {

	return func(addr string) (session.Conn, error) {

		conn, err := net.Dial(network, addr)
		if err != nil {
			return nil, err
		}

		return &Conn{conn: conn}, nil
	}
}
