package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Conn struct {
	conn *websocket.Conn
	addr string
	data *sync.Map
}

func (conn *Conn) Addr() string {
	return conn.addr
}

func (conn *Conn) Deadline(time time.Time) error {
	conn.conn.SetReadDeadline(time)
	conn.conn.SetWriteDeadline(time)
	return nil
}

func (conn *Conn) Read() ([]byte, error) {
	_, raw, err := conn.conn.ReadMessage()
	return raw, err
}

func (conn *Conn) Write(stream []byte) error {
	return conn.conn.WriteMessage(websocket.BinaryMessage, stream)
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

func (conn *Conn) Data() *sync.Map {
	return conn.data
}
