package session

import "time"

type Conn interface {
	Addr() string
	Deadline(time.Time) error
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}
