package websocket

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/laconiz/metis/nets/acceptor"
	"github.com/laconiz/metis/nets/queue"
	"github.com/laconiz/metis/nets/session"
	"net/http"
)

func Acceptor(limit int64) acceptor.Dialer {

	return func(addr string) (acceptor.Listener, error) {

		listener := &Listener{
			limit: limit,
			queue: queue.New(capacity),
		}

		engine := gin.New()
		engine.Use(gin.Recovery())
		engine.Any(Path, listener.upgrade)

		listener.server = &http.Server{
			Addr:    addr,
			Handler: engine,
		}
		go listener.run()

		return listener, nil
	}
}

// 侦听器
type Listener struct {
	limit  int64        // 读取限制
	queue  *queue.Queue // 会话队列
	server *http.Server // HTTP服务
}

// 获取会话
func (listener *Listener) Accept() (session.Conn, error) {

	// 从会话队列取出会话
	conn, closed := listener.queue.Pop()
	if closed {
		return nil, errClosed
	}

	return conn.(session.Conn), nil
}

// 运行服务
func (listener *Listener) run() {
	if err := listener.server.ListenAndServe(); err != nil {
		listener.queue.Close()
	}
}

// 停止服务
func (listener *Listener) Close() error {
	return listener.server.Shutdown(context.Background())
}

// 升级会话
func (listener *Listener) upgrade(context *gin.Context) {

	request := context.Request

	// 升级会话
	conn, err := upgrader.Upgrade(context.Writer, request, nil)
	if err != nil {
		return
	}

	// 设置读取限制
	conn.SetReadLimit(listener.limit)

	// 代理转发地址
	addr := request.RemoteAddr
	if addr == "" {
		addr = conn.RemoteAddr().String()
	}

	// 插入会话队列
	listener.queue.Push(&Conn{conn: conn, addr: addr})
}

// 升级器
var upgrader = &websocket.Upgrader{
	EnableCompression: true,
	HandshakeTimeout:  HandshakeTimeout,
}

// 侦听器关闭
var errClosed = errors.New(acceptor.StrClosed)

// 会话队列限制
var capacity = 64

func init() {
	gin.SetMode(gin.ReleaseMode)
}
