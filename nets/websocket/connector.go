package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/laconiz/metis/nets/connector"
	"github.com/laconiz/metis/nets/session"
	"net/http"
)

func Connector() connector.Dialer {

	return func(addr string) (session.Conn, error) {

		// 创建连接
		conn, _, err := dialer.Dial(addr, header)
		if err != nil {
			return nil, err
		}

		address := conn.RemoteAddr().String()
		return &Conn{conn: conn, addr: address}, nil
	}
}

// 连接器
var dialer = &websocket.Dialer{
	HandshakeTimeout: HandshakeTimeout,
}

// 连接头
var header = http.Header{}
