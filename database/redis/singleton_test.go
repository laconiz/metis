package redis

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSingleton_Exec(t *testing.T) {

	assert(client.Key().Delete(KeyA) == nil)

	var count int32
	var success int32
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {

		wg.Add(1)

		go func() {

			if ok, _ := client.Singleton(KeyA).Exec(func() {
				atomic.AddInt32(&count, 1)
			}); ok {
				atomic.AddInt32(&success, 1)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	assert(count == 1, success == 1)

	assert(client.Key().Delete(KeyA) == nil)
}
