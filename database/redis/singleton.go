package redis

import (
	"time"
)

type Singleton struct {
	key    string
	client *Redis
}

func (singleton *Singleton) Exec(handler func()) (bool, error) {
	key := singleton.client.Key()
	ok, err := key.SetNX(singleton.key, time.Now().String())
	if ok {
		handler()
	}
	return ok, err
}
