package redis

import (
	"encoding/hex"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
	"time"
)

var (
	ErrAtomicLockFailed   = errors.New("lock failed")
	ErrAtomicUnlockFailed = errors.New("unlock failed")
)

type Atomic struct {
	key     string
	expired int64
	ticker  time.Duration
	timeout time.Duration
	value   string
	client  *Client
}

func (atomic *Atomic) Expired(duration time.Duration) *Atomic {
	atomic.expired = int64(duration / time.Second)
	return atomic
}

func (atomic *Atomic) Ticker(duration time.Duration) *Atomic {
	atomic.ticker = duration
	return atomic
}

func (atomic *Atomic) Timeout(duration time.Duration) *Atomic {
	atomic.timeout = duration
	return atomic
}

func (atomic *Atomic) Lock() error {

	atomic.value = hex.EncodeToString(uuid.NewV1().Bytes())

	ticker := time.NewTicker(atomic.ticker)
	defer ticker.Stop()

	deadline := time.Now().Add(atomic.timeout)

	key := atomic.client.Key()

	for {

		ok, err := key.SetNEX(atomic.key, atomic.value, atomic.expired)
		if err != nil {
			return err
		} else if ok {
			return nil
		}

		if deadline.Before(time.Now()) {
			return ErrAtomicLockFailed
		}

		<-ticker.C
	}
}

func (atomic *Atomic) Unlock() error {

	script := atomic.client.Eval(scriptAtomicUnlock)
	reply, err := script.Exec(atomic.key, atomic.value)

	if ok, err := redis.Bool(reply, err); err != nil {
		return err
	} else if !ok {
		return ErrAtomicUnlockFailed
	} else {
		return nil
	}
}

func (atomic *Atomic) Exec(handler func()) (ok bool, err error) {

	if err := atomic.Lock(); err != nil {
		return false, err
	}

	defer func() { err = atomic.Unlock() }()

	defer func() {
		if err := recover(); err != nil {
			if logger := atomic.client.option.Logger; logger != nil {
				logger.Data(err).Error("execute error")
			}
		}
	}()

	handler()
	ok = true
	return
}

var scriptAtomicUnlock = &Script{Name: "AtomicUnlock", Script: redis.NewScript(1, luaAtomicUnlock)}

var luaAtomicUnlock = `
	if redis.call('GET', KEYS[1]) == ARGV[1] then
		return redis.call('DEL', KEYS[1])
	end
	return 0
`
