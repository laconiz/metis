package socket

import (
	"github.com/laconiz/metis/nets/acceptor"
	"github.com/laconiz/metis/nets/session"
	"net"
)

func Acceptor(limit int32) acceptor.Dialer {

	return func(addr string) (acceptor.Listener, error) {

		listener, err := net.Listen(network, addr)
		if err != nil {
			return nil, err
		}

		return &Listener{limit: limit, listener: listener}, nil
	}
}

// 侦听器
type Listener struct {
	limit    int32        // 读取限制
	listener net.Listener // TCP服务
}

// 获取会话
func (listener *Listener) Accept() (session.Conn, error) {

	// 建立连接
	conn, err := listener.listener.Accept()
	if err != nil {
		return nil, err
	}

	return &Conn{conn: conn, limit: listener.limit}, nil
}

func (listener *Listener) Close() error {
	return listener.listener.Close()
}
