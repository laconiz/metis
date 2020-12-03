package socket

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

type Conn struct {
	conn  net.Conn
	limit int32
}

func (conn *Conn) Addr() string {
	return conn.conn.RemoteAddr().String()
}

func (conn *Conn) Deadline(time time.Time) error {
	return conn.conn.SetDeadline(time)
}

func (conn *Conn) Read() ([]byte, error) {

	var size int32
	if err := binary.Read(conn.conn, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("read size error: %w", err)
	}

	if size <= 0 || (conn.limit > 0 && size > conn.limit) {
		return nil, fmt.Errorf("invalid size: %d", size)
	}

	stream := make([]byte, size)
	if _, err := conn.conn.Read(stream); err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	return stream, nil
}

func (conn *Conn) Write(stream []byte) error {

	if len(stream) == 0 {
		return fmt.Errorf("nil stream")
	}

	size := int32(len(stream))
	if err := binary.Write(conn.conn, binary.LittleEndian, size); err != nil {
		return fmt.Errorf("write size error: %w", err)
	}

	if _, err := conn.conn.Write(stream); err != nil {
		return fmt.Errorf("write body error: %w", err)
	}

	return nil
}

func (conn *Conn) Close() error {
	return conn.conn.Close()
}

func (conn *Conn) Data() *sync.Map {
	return &sync.Map{}
}
