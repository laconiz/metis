package session

import (
	"sync"
	"time"
)

type Conn interface {
	Addr() string
	Deadline(time.Time) error
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
	Data() *sync.Map
}
