package redis

import (
	"sync"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	executed, err := client.Atomic(KeyA).Exec(func() {})
	assert(err == nil, executed)

	executed, err = client.Atomic(KeyA).
		Expired(time.Second).
		Timeout(time.Second * 2).
		Ticker(time.Millisecond * 100).
		Exec(func() { time.Sleep(time.Millisecond * 1100) })
	assert(err == ErrAtomicUnlockFailed, executed)

	assert(client.Key().Set(KeyA, ValueA) == nil)
	executed, err = client.Atomic(KeyA).Exec(func() {})
	assert(err == ErrAtomicLockFailed, !executed)

	assert(client.Key().Delete(KeyA) == nil)

	var count int32
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			executed, err := client.Atomic(KeyA).Exec(func() {
				count++
				time.Sleep(time.Millisecond * 100)
			})
			assert(err == nil, executed)
			wg.Done()
		}()
	}

	wg.Wait()
	assert(count == 20)

	assert(client.Key().Delete(KeyA) == nil)
}
